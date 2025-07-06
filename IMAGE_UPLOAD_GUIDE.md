# Image Upload Feature Test Guide

## Overview
The Agricultural Equipment Store API now supports multiple image uploads for products. You can:
- Upload local image files (JPEG, PNG, GIF, WebP)
- Provide image URLs
- Mix both uploaded files and URLs in one request
- Support multiple images per product

## API Endpoints

### Create Product with Images
- **URL**: `POST /api/products`
- **Content-Type**: `multipart/form-data` (for file uploads) or `application/json` (for URLs only)
- **Authentication**: Bearer token required (admin only)

#### Form Fields (multipart/form-data):
- `name` (required): Product name
- `description`: Product description  
- `price` (required): Product price
- `category` (required): Product category
- `brand`: Product brand
- `stock` (required): Stock quantity
- `image_urls`: Comma-separated image URLs
- `images`: Multiple image files

#### JSON Fields (application/json):
```json
{
  "name": "Product Name",
  "description": "Product Description",
  "price": 99.99,
  "category": "Category",
  "brand": "Brand",
  "stock": 10,
  "image_url": "https://example.com/image.jpg",
  "image_urls": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"]
}
```

### Update Product with Images
- **URL**: `PUT /api/products/{id}`
- **Content-Type**: `multipart/form-data` (for file uploads) or `application/json` (for URLs only)
- **Authentication**: Bearer token required (admin only)

Same fields as create, plus:
- `is_active`: Boolean for product status

## Testing

### 1. Using the HTML Test Page
1. Open `test_image_upload.html` in your browser
2. Click "Login as Admin" to authenticate
3. Fill in the product form
4. Select image files and/or add image URLs
5. Submit the form

### 2. Using curl (Command Line)

#### Login to get token:
```bash
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@agricultural.com", "password": "password123"}'
```

#### Create product with file upload:
```bash
curl -X POST http://localhost:8082/api/products \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "name=Test Product" \
  -F "description=Test Description" \
  -F "price=99.99" \
  -F "category=Tractors" \
  -F "brand=Test Brand" \
  -F "stock=10" \
  -F "images=@/path/to/image1.jpg" \
  -F "images=@/path/to/image2.jpg" \
  -F "image_urls=https://example.com/image3.jpg"
```

#### Create product with URLs only:
```bash
curl -X POST http://localhost:8082/api/products \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "URL Product",
    "description": "Product with URL images",
    "price": 149.99,
    "category": "Plows",
    "brand": "Brand",
    "stock": 5,
    "image_urls": [
      "https://example.com/plow1.jpg",
      "https://example.com/plow2.jpg"
    ]
  }'
```

### 3. View Uploaded Images
- Uploaded images are served at: `http://localhost:8082/uploads/products/{filename}`
- Check the product response for the full image URLs

## File Upload Constraints
- **Max file size**: 5MB per file
- **Allowed types**: JPEG, PNG, GIF, WebP
- **Max files**: No limit set (reasonable use expected)
- **Storage**: Files saved to `uploads/products/` directory

## Image Structure
Each product now contains an `images` array with objects containing:
- `id`: Unique identifier
- `url`: Full URL to access the image
- `filename`: Original filename (for uploaded files)
- `file_path`: Server path (for uploaded files)
- `file_size`: File size in bytes
- `mime_type`: MIME type
- `is_url`: Boolean indicating if it's a URL or uploaded file
- `is_primary`: Boolean indicating the main product image
- `created_at`: Timestamp

## Backward Compatibility
- The legacy `image_url` field is still supported
- Existing products with `image_url` will continue to work
- New products can use both old and new image fields

## Error Handling
- Invalid file types return 400 error
- File size exceeded returns 400 error
- Missing authentication returns 401 error
- Non-admin users return 403 error
- Upload failures clean up any partially uploaded files

## API Documentation
- Full API documentation available at: `http://localhost:8082/swagger/index.html`
- Interactive testing available through Swagger UI
