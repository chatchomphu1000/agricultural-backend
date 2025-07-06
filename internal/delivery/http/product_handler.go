package http

import (
	"agricultural-equipment-store/internal/domain"
	"agricultural-equipment-store/internal/usecase"
	"agricultural-equipment-store/internal/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	uploadConfig   *utils.UploadConfig
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		uploadConfig:   utils.NewUploadConfig(),
	}
}

// CreateProduct handles creating a new product
// @Summary Create a new product
// @Description Create a new product (admin only). Supports both JSON with image URLs and multipart form with file uploads.
// @Tags products
// @Accept json,multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateProductRequest true "Product creation request (JSON)"
// @Param name formData string true "Product name (Form)"
// @Param description formData string false "Product description (Form)"
// @Param price formData number true "Product price (Form)"
// @Param category formData string true "Product category (Form)"
// @Param brand formData string false "Product brand (Form)"
// @Param stock formData integer true "Product stock (Form)"
// @Param image_urls formData string false "Comma-separated image URLs (Form)"
// @Param images formData file false "Product images (Form, multiple files allowed)"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		h.createProductWithFiles(c)
	} else {
		h.createProductWithJSON(c)
	}
}

// createProductWithJSON handles JSON-based product creation
func (h *ProductHandler) createProductWithJSON(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUseCase.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// createProductWithFiles handles multipart form-based product creation with file uploads
func (h *ProductHandler) createProductWithFiles(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Extract basic product data
	req := domain.CreateProductRequest{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Category:    c.PostForm("category"),
		Brand:       c.PostForm("brand"),
	}

	// Parse price
	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		} else {
			req.Price = price
		}
	}

	// Parse stock
	if stockStr := c.PostForm("stock"); stockStr != "" {
		if stock, err := strconv.Atoi(stockStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock format"})
			return
		} else {
			req.Stock = stock
		}
	}

	// Validate required fields
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
		return
	}
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product price must be greater than 0"})
		return
	}
	if req.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product category is required"})
		return
	}

	// Handle image URLs (comma-separated)
	if imageURLs := c.PostForm("image_urls"); imageURLs != "" {
		req.ImageURLs = strings.Split(imageURLs, ",")
		for i, url := range req.ImageURLs {
			req.ImageURLs[i] = strings.TrimSpace(url)
		}
	}

	// Handle legacy single image URL
	if imageURL := c.PostForm("image_url"); imageURL != "" {
		req.ImageURL = strings.TrimSpace(imageURL)
	}

	// Handle file uploads
	var uploadedImages []domain.ProductImage
	if form := c.Request.MultipartForm; form != nil && form.File["images"] != nil {
		for _, fileHeader := range form.File["images"] {
			result, err := h.uploadConfig.SaveFile(fileHeader)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to upload file %s: %v", fileHeader.Filename, err)})
				return
			}

			// Generate image URL for serving
			baseURL := fmt.Sprintf("%s://%s", c.Request.URL.Scheme, c.Request.Host)
			if baseURL == "://" {
				baseURL = "http://localhost:8082" // fallback for local development
			}
			imageURL := utils.GenerateImageURL(result.FilePath, baseURL)

			uploadedImages = append(uploadedImages, domain.ProductImage{
				ID:        result.ID,
				URL:       imageURL,
				Filename:  result.Filename,
				FilePath:  result.FilePath,
				FileSize:  result.FileSize,
				MimeType:  result.MimeType,
				IsURL:     false,
				IsPrimary: len(uploadedImages) == 0, // First image is primary
				CreatedAt: time.Now(),
			})
		}
	}

	// Create product with enhanced request
	product, err := h.productUseCase.CreateProductWithImages(c.Request.Context(), req, uploadedImages)
	if err != nil {
		// Clean up uploaded files on error
		for _, img := range uploadedImages {
			if !img.IsURL {
				h.uploadConfig.DeleteFile(img.FilePath)
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct handles getting a product by ID
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.productUseCase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts handles getting products with filtering and pagination
// @Summary Get products
// @Description Get products with optional filtering and pagination
// @Tags products
// @Produce json
// @Param category query string false "Category filter"
// @Param brand query string false "Brand filter"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param search query string false "Search in name and description"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse query parameters
	filter := domain.ProductFilter{}

	filter.Category = c.Query("category")
	filter.Brand = c.Query("brand")
	filter.Search = c.Query("search")

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = maxPrice
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	// Set default active filter to true for public endpoint
	isActive := true
	filter.IsActive = &isActive

	products, count, err := h.productUseCase.GetProducts(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pagination info
	var totalPages int64
	if filter.Limit > 0 {
		totalPages = (count + int64(filter.Limit) - 1) / int64(filter.Limit)
	} else {
		totalPages = 1
	}

	response := gin.H{
		"products":    products,
		"total":       count,
		"page":        filter.Page,
		"limit":       filter.Limit,
		"total_pages": totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProduct handles updating a product
// @Summary Update product
// @Description Update a product (admin only). Supports both JSON with image URLs and multipart form with file uploads.
// @Tags products
// @Accept json,multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body domain.UpdateProductRequest true "Product update request (JSON)"
// @Param name formData string false "Product name (Form)"
// @Param description formData string false "Product description (Form)"
// @Param price formData number false "Product price (Form)"
// @Param category formData string false "Product category (Form)"
// @Param brand formData string false "Product brand (Form)"
// @Param stock formData integer false "Product stock (Form)"
// @Param is_active formData boolean false "Product active status (Form)"
// @Param image_urls formData string false "Comma-separated image URLs (Form)"
// @Param images formData file false "Product images (Form, multiple files allowed)"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		h.updateProductWithFiles(c)
	} else {
		h.updateProductWithJSON(c)
	}
}

// updateProductWithJSON handles JSON-based product updates
func (h *ProductHandler) updateProductWithJSON(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUseCase.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// updateProductWithFiles handles multipart form-based product updates with file uploads
func (h *ProductHandler) updateProductWithFiles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Parse multipart form
	err = c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Extract product data
	req := domain.UpdateProductRequest{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Category:    c.PostForm("category"),
		Brand:       c.PostForm("brand"),
	}

	// Parse price
	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		} else {
			req.Price = price
		}
	}

	// Parse stock
	if stockStr := c.PostForm("stock"); stockStr != "" {
		if stock, err := strconv.Atoi(stockStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock format"})
			return
		} else {
			req.Stock = stock
		}
	}

	// Parse is_active
	if isActiveStr := c.PostForm("is_active"); isActiveStr != "" {
		if isActiveStr == "true" || isActiveStr == "1" {
			isActive := true
			req.IsActive = &isActive
		} else if isActiveStr == "false" || isActiveStr == "0" {
			isActive := false
			req.IsActive = &isActive
		}
	}

	// Handle image URLs (comma-separated)
	if imageURLs := c.PostForm("image_urls"); imageURLs != "" {
		req.ImageURLs = strings.Split(imageURLs, ",")
		for i, url := range req.ImageURLs {
			req.ImageURLs[i] = strings.TrimSpace(url)
		}
	}

	// Handle legacy single image URL
	if imageURL := c.PostForm("image_url"); imageURL != "" {
		req.ImageURL = strings.TrimSpace(imageURL)
	}

	// Handle file uploads
	var uploadedImages []domain.ProductImage
	if form := c.Request.MultipartForm; form != nil && form.File["images"] != nil {
		for _, fileHeader := range form.File["images"] {
			result, err := h.uploadConfig.SaveFile(fileHeader)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to upload file %s: %v", fileHeader.Filename, err)})
				return
			}

			// Generate image URL for serving
			baseURL := fmt.Sprintf("%s://%s", c.Request.URL.Scheme, c.Request.Host)
			if baseURL == "://" {
				baseURL = "http://localhost:8082" // fallback for local development
			}
			imageURL := utils.GenerateImageURL(result.FilePath, baseURL)

			uploadedImages = append(uploadedImages, domain.ProductImage{
				ID:        result.ID,
				URL:       imageURL,
				Filename:  result.Filename,
				FilePath:  result.FilePath,
				FileSize:  result.FileSize,
				MimeType:  result.MimeType,
				IsURL:     false,
				IsPrimary: len(uploadedImages) == 0, // First image is primary
				CreatedAt: time.Now(),
			})
		}
	}

	// Update product with enhanced request
	product, err := h.productUseCase.UpdateProductWithImages(c.Request.Context(), id, req, uploadedImages)
	if err != nil {
		// Clean up uploaded files on error
		for _, img := range uploadedImages {
			if !img.IsURL {
				h.uploadConfig.DeleteFile(img.FilePath)
			}
		}
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handles deleting a product
// @Summary Delete product
// @Description Delete a product (admin only)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	err = h.productUseCase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
