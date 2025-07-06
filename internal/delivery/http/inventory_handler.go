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

// InventoryHandler handles inventory related endpoints
type InventoryHandler struct {
	inventoryUseCase *usecase.InventoryUseCase
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryUseCase *usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		inventoryUseCase: inventoryUseCase,
	}
}

// UpdateStock handles updating product stock
// @Summary Update product stock
// @Description Update stock quantity for a specific product
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body domain.StockUpdateRequest true "Stock update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /inventories/{id}/stock [put]
func (h *InventoryHandler) UpdateStock(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req domain.StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.inventoryUseCase.UpdateStock(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

// GetLowStockProducts handles getting products with low stock
// @Summary Get low stock products
// @Description Get products with stock below threshold
// @Tags inventory
// @Produce json
// @Security BearerAuth
// @Param threshold query int false "Stock threshold (default: 10)"
// @Success 200 {array} domain.LowStockProduct
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventories/low-stock [get]
func (h *InventoryHandler) GetLowStockProducts(c *gin.Context) {
	threshold := 10 // Default threshold
	if thresholdStr := c.Query("threshold"); thresholdStr != "" {
		if t, err := strconv.Atoi(thresholdStr); err == nil && t > 0 {
			threshold = t
		}
	}

	products, err := h.inventoryUseCase.GetLowStockProducts(c.Request.Context(), threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetStockSummary handles getting stock summary
// @Summary Get stock summary
// @Description Get overall stock summary including totals and category breakdown
// @Tags inventory
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.StockSummary
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventories/summary [get]
func (h *InventoryHandler) GetStockSummary(c *gin.Context) {
	summary, err := h.inventoryUseCase.GetStockSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// SaleHandler handles sales related endpoints
type SaleHandler struct {
	saleUseCase *usecase.SaleUseCase
}

// NewSaleHandler creates a new sale handler
func NewSaleHandler(saleUseCase *usecase.SaleUseCase) *SaleHandler {
	return &SaleHandler{
		saleUseCase: saleUseCase,
	}
}

// CreateSale handles creating a new sale
// @Summary Create a new sale
// @Description Create a new sale transaction
// @Tags sales
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateSaleRequest true "Sale creation request"
// @Success 201 {object} domain.Sale
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sales [post]
func (h *SaleHandler) CreateSale(c *gin.Context) {
	var req domain.CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sale, err := h.saleUseCase.CreateSale(c.Request.Context(), req)
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

// GetSalesSummary handles getting sales summary
// @Summary Get sales summary
// @Description Get sales summary for a specific period
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} domain.SalesSummary
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sales/summary [get]
func (h *SaleHandler) GetSalesSummary(c *gin.Context) {
	var fromDate, toDate time.Time
	var err error

	if fromStr := c.Query("from"); fromStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format (use YYYY-MM-DD)"})
			return
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		toDate, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format (use YYYY-MM-DD)"})
			return
		}
		// Set to end of day
		toDate = toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	summary, err := h.saleUseCase.GetSalesSummary(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetSales handles getting sales with filtering
// @Summary Get sales
// @Description Get sales with optional filtering by date range
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param product_id query string false "Product ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} domain.Sale
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sales [get]
func (h *SaleHandler) GetSales(c *gin.Context) {
	var filter domain.SaleFilter
	var err error

	if fromStr := c.Query("from"); fromStr != "" {
		filter.FromDate, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format (use YYYY-MM-DD)"})
			return
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		filter.ToDate, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format (use YYYY-MM-DD)"})
			return
		}
		// Set to end of day
		filter.ToDate = filter.ToDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	if productIDStr := c.Query("product_id"); productIDStr != "" {
		filter.ProductID, err = primitive.ObjectIDFromHex(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
			return
		}
	}

	// Pagination
	filter.Page = 1
	filter.Limit = 10
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	sales, err := h.saleUseCase.GetSalesByFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sales)
}

// GetSalesByProduct handles getting sales grouped by product
// @Summary Get sales by product
// @Description Get sales data grouped by product for analytics
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} domain.ProductSales
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sales/by-product [get]
func (h *SaleHandler) GetSalesByProduct(c *gin.Context) {
	var fromDate, toDate time.Time
	var err error

	if fromStr := c.Query("from"); fromStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format (use YYYY-MM-DD)"})
			return
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		toDate, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format (use YYYY-MM-DD)"})
			return
		}
		// Set to end of day
		toDate = toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	productSales, err := h.saleUseCase.GetSalesByProduct(c.Request.Context(), fromDate, toDate)
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
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {string} string "CSV data"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sales/export [get]
func (h *SaleHandler) ExportSales(c *gin.Context) {
	var fromDate, toDate time.Time
	var err error

	if fromStr := c.Query("from"); fromStr != "" {
		fromDate, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format (use YYYY-MM-DD)"})
			return
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		toDate, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format (use YYYY-MM-DD)"})
			return
		}
		// Set to end of day
		toDate = toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	sales, err := h.saleUseCase.GetSalesByDateRange(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate CSV
	csvData := "ID,Product ID,Quantity,Price,Total,Date Sold\n"
	for _, sale := range sales {
		csvData += sale.ID.Hex() + "," +
			sale.ProductID.Hex() + "," +
			strconv.Itoa(sale.Quantity) + "," +
			strconv.FormatFloat(sale.Price, 'f', 2, 64) + "," +
			strconv.FormatFloat(sale.Total, 'f', 2, 64) + "," +
			sale.DateSold.Format("2006-01-02 15:04:05") + "\n"
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=sales_export.csv")
	c.String(http.StatusOK, csvData)
}
