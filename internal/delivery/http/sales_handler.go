package http

import (
	"agricultural-equipment-store/internal/domain"
	"agricultural-equipment-store/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SalesHandler handles sales endpoints
type SalesHandler struct {
	salesUseCase *usecase.SalesUseCase
}

// NewSalesHandler creates a new sales handler
func NewSalesHandler(salesUseCase *usecase.SalesUseCase) *SalesHandler {
	return &SalesHandler{
		salesUseCase: salesUseCase,
	}
}

// CreateSale handles creating a new sale
// @Summary Create a new sale
// @Description Create a new sale record
// @Tags sales
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateSaleRequest true "Sale creation request"
// @Success 201 {object} domain.Sale
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sales [post]
func (h *SalesHandler) CreateSale(c *gin.Context) {
	var req domain.CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sale, err := h.salesUseCase.CreateSale(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sale)
}

// GetSales handles getting sales with filtering
// @Summary Get sales
// @Description Get sales list with optional filtering
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param product_id query string false "Product ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} domain.Sale
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sales [get]
func (h *SalesHandler) GetSales(c *gin.Context) {
	// Parse query parameters
	var filter domain.SaleFilter

	// Parse dates
	if fromStr := c.Query("from"); fromStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromStr); err == nil {
			filter.FromDate = fromDate
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if toDate, err := time.Parse("2006-01-02", toStr); err == nil {
			filter.ToDate = toDate
		}
	}

	// Parse product ID
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if productID, err := primitive.ObjectIDFromHex(productIDStr); err == nil {
			filter.ProductID = productID
		}
	}

	// Parse pagination
	if pageStr := c.DefaultQuery("page", "1"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}

	if limitStr := c.DefaultQuery("limit", "10"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// Get sales
	sales, err := h.salesUseCase.GetSales(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sales)
}

// GetSalesSummary handles getting sales summary
// @Summary Get sales summary
// @Description Get sales summary for current month or specified period
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} domain.SalesSummary
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sales/summary [get]
func (h *SalesHandler) GetSalesSummary(c *gin.Context) {
	// Parse dates or use current month
	var fromDate, toDate time.Time

	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			fromDate = parsed
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			toDate = parsed
		}
	}

	// If no dates provided, use current month
	if fromDate.IsZero() && toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		toDate = fromDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Get sales summary
	summary, err := h.salesUseCase.GetSalesSummary(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetSalesByProduct handles getting sales by product
// @Summary Get sales by product
// @Description Get sales data grouped by product
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Success 200 {array} domain.ProductSales
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sales/by-product [get]
func (h *SalesHandler) GetSalesByProduct(c *gin.Context) {
	// Parse dates or use current month
	var fromDate, toDate time.Time

	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			fromDate = parsed
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			toDate = parsed
		}
	}

	// If no dates provided, use current month
	if fromDate.IsZero() && toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		toDate = fromDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Get sales by product
	productSales, err := h.salesUseCase.GetSalesByProduct(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productSales)
}

// ExportSales handles exporting sales data
// @Summary Export sales data
// @Description Export sales data as CSV
// @Tags sales
// @Produce text/csv
// @Security BearerAuth
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param format query string false "Export format (csv)" default(csv)
// @Success 200 {string} string "CSV data"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sales/export [get]
func (h *SalesHandler) ExportSales(c *gin.Context) {
	// Parse dates or use current month
	var fromDate, toDate time.Time

	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			fromDate = parsed
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			toDate = parsed
		}
	}

	// If no dates provided, use current month
	if fromDate.IsZero() && toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		toDate = fromDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Get format
	format := c.DefaultQuery("format", "csv")
	if format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only CSV format is supported"})
		return
	}

	// Export sales
	csvData, err := h.salesUseCase.ExportSales(c.Request.Context(), fromDate, toDate, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for CSV download
	filename := "sales_" + fromDate.Format("2006-01-02") + "_to_" + toDate.Format("2006-01-02") + ".csv"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv")
	c.String(http.StatusOK, csvData)
}
