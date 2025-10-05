package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/warleon/ms4-compliance-service/internal/dto"
	"github.com/warleon/ms4-compliance-service/internal/repository/rules"
	"github.com/warleon/ms4-compliance-service/internal/service"
)

type ComplianceHandler struct {
	service *service.ComplianceService
}

func NewComplianceHandler(s *service.ComplianceService) *ComplianceHandler {
	return &ComplianceHandler{service: s}
}

// ValidateTransaction godoc
// @Summary Validate a transaction
// @Description Validates a transaction against compliance rules and returns a decision
// @Tags compliance
// @Accept json
// @Produce json
// @Param transaction body dto.Transaction true "Transaction data"
// @Success 200 {object} rules.Decision
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/validateTransaction [post]
func (h *ComplianceHandler) ValidateTransaction(c *gin.Context) {
	var tx dto.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dec, err := h.service.ValidateTransaction(c.Request.Context(), tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dec)
}

// CreateRule godoc
// @Summary Create a new rule
// @Description Creates a compliance rule in the system
// @Tags rules
// @Accept json
// @Produce json
// @Param rule body rules.Rule true "Rule data"
// @Success 201 {object} rules.Rule
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/rules [post]
func (h *ComplianceHandler) CreateRule(c *gin.Context) {
	var r rules.Rule
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateRule(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, r)
}

// ListRules godoc
// @Summary List all rules
// @Description Retrieves a paginated list of compliance rules
// @Tags rules
// @Accept json
// @Produce json
// @Param size query int false "Number of results per page (default 50)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {array} rules.Rule
// @Failure 500 {object} map[string]string
// @Router /api/v1/rules [get]
func (h *ComplianceHandler) ListRules(c *gin.Context) {
	size := 50
	offset := 0
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			size = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}
	rs, err := h.service.ListRules(c.Request.Context(), size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rs)
}

// GetRule godoc
// @Summary Get rule by ID
// @Description Retrieves a specific compliance rule by its ID
// @Tags rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} rules.Rule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/rules/{id} [get]
func (h *ComplianceHandler) GetRule(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	r, err := h.service.GetRule(c.Request.Context(), uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}
	c.JSON(http.StatusOK, r)
}

// UpdateRule godoc
// @Summary Update an existing rule
// @Description Updates a compliance rule by ID
// @Tags rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Param rule body rules.Rule true "Updated rule data"
// @Success 200 {object} rules.Rule
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/rules/{id} [put]
func (h *ComplianceHandler) UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var r rules.Rule
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.ID = uint(id64)
	if err := h.service.UpdateRule(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

// DeleteRule godoc
// @Summary Delete a rule
// @Description Deletes a compliance rule by ID
// @Tags rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/rules/{id} [delete]
func (h *ComplianceHandler) DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteRule(c.Request.Context(), uint(id64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
