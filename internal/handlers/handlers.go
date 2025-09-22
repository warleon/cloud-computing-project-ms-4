package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/warleon/ms4-compliance-service/internal/repository"
	"github.com/warleon/ms4-compliance-service/internal/service"
)

type ComplianceHandler struct {
	service *service.ComplianceService
}

func NewComplianceHandler(s *service.ComplianceService) *ComplianceHandler {
	return &ComplianceHandler{service: s}
}

func (h *ComplianceHandler) ValidateTransaction(c *gin.Context) {
	var in service.ValidateTransactionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.ValidateTransaction(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ComplianceHandler) GetRiskScore(c *gin.Context) {
	customerID := c.Param("customerId")
	score, err := h.service.GetRiskScore(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"customerId": customerID, "riskScore": score})
}

func (h *ComplianceHandler) CreateRule(c *gin.Context) {
	var r repository.Rule
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Repo.CreateRule(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, r)
}

func (h *ComplianceHandler) ListRules(c *gin.Context) {
	rules, err := h.service.Repo.ListRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}
