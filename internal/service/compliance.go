package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/warleon/ms4-compliance-service/internal/repository"
)

// ComplianceService contains business logic.
type ComplianceService struct {
	Repo  repository.Repository
	Fraud FraudClient
}

func NewComplianceService(repo repository.Repository, fraud FraudClient) *ComplianceService {
	return &ComplianceService{Repo: repo, Fraud: fraud}
}

// ValidateTransactionInput DTO
type ValidateTransactionInput struct {
	TransactionID string  `json:"transactionId" binding:"required"`
	CustomerID    string  `json:"customerId" binding:"required"`
	FromAccount   string  `json:"fromAccount" binding:"required"`
	ToAccount     string  `json:"toAccount" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required"`
	Metadata      string  `json:"metadata"`
}

// ValidateTransactionResult DTO
type ValidateTransactionResult struct {
	Decision string  `json:"decision"`
	Reason   string  `json:"reason"`
	Score    float64 `json:"score"`
}

func (s *ComplianceService) ValidateTransaction(ctx context.Context, in ValidateTransactionInput) (*ValidateTransactionResult, error) {
	logrus.WithFields(logrus.Fields{"tx": in.TransactionID, "customer": in.CustomerID}).Info("starting validation")

	// 1. check rules type: amount threshold
	rules, err := s.Repo.FindRulesByType("amount_threshold")
	if err != nil {
		return nil, err
	}

	// compute rule-based score
	score := 0.0
	for _, r := range rules {
		// try to parse Params as JSON: {"threshold":10000, "weight":0.7}
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(r.Params), &params); err != nil {
			continue
		}
		threshold := toFloat(params["threshold"])
		weight := toFloat(params["weight"])
		if in.Amount >= threshold {
			score += weight
		}
	}

	// 2. sanctions list: simple rule lookup (params could be a CSV of blacklisted accounts)
	sanctions, _ := s.Repo.FindRulesByType("sanctions_list")
	for _, r := range sanctions {
		if strings.Contains(r.Params, in.ToAccount) || strings.Contains(r.Params, in.FromAccount) || strings.Contains(r.Params, in.CustomerID) {
			score += 1.0
		}
	}

	// 3. consult external fraud API
	fraudPayload := map[string]interface{}{
		"transactionId": in.TransactionID,
		"amount":        in.Amount,
		"from":          in.FromAccount,
		"to":            in.ToAccount,
		"customerId":    in.CustomerID,
	}
	fraudResp, err := s.Fraud.Evaluate(ctx, fraudPayload)
	if err == nil {
		if v, ok := fraudResp["score"].(float64); ok {
			score += v
		}
	}

	// normalize
	if score > 1.0 {
		score = score / 2.0
	} // crude normalization

	decision := "approve"
	reason := ""
	if score >= 1.5 {
		decision = "reject"
		reason = "high risk based on rules and external signals"
	} else if score >= 0.7 {
		decision = "review"
		reason = "requires manual review"
	}

	// 4. persist audit
	audit := &repository.AuditLog{
		TransactionID: in.TransactionID,
		CustomerID:    in.CustomerID,
		Decision:      decision,
		Reason:        reason,
		Meta:          fmt.Sprintf("score=%.2f", score),
		CreatedAt:     time.Now(),
	}
	s.Repo.CreateAudit(audit)

	res := &ValidateTransactionResult{Decision: decision, Reason: reason, Score: score}
	return res, nil
}

func (s *ComplianceService) GetRiskScore(ctx context.Context, customerID string) (float64, error) {
	// simple aggregator: look at recent audit logs and compute a risk score
	audits, err := s.Repo.ListAuditsForCustomer(customerID, 50)
	if err != nil {
		return 0, err
	}
	var score float64
	for _, a := range audits {
		if a.Decision == "reject" {
			score += 1.0
		}
		if a.Decision == "review" {
			score += 0.5
		}
	}
	// normalize into 0..1
	if score > 5 {
		score = 5
	}
	return score / 5.0, nil
}

func toFloat(v interface{}) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	default:
		return 0
	}
}
