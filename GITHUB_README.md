# Agricultural Equipment Store API

A comprehensive REST API for managing agricultural equipment store operations, built with Go and MongoDB.

## 🚀 Features

### 🔐 Authentication & Authorization
- JWT-based authentication
- Role-based access control (Admin/User)
- Secure password hashing
- Token-based API access

### 📦 Product Management
- **Multiple Image Support** - Upload local files + image URLs
- CRUD operations for products
- Advanced filtering and search
- Category management
- Stock tracking
- Price management

### 🖼️ Image Upload System
- **Local file uploads** with validation
- **Image URLs** support
- **Multiple images** per product
- File type validation (JPEG, PNG, GIF, WebP)
- File size limits (5MB max)
- Automatic image serving
- Backward compatibility with legacy `image_url` field

### 📊 Inventory Management
- Stock level monitoring
- Low stock alerts
- Inventory summaries
- Stock value calculations
- Category-based reporting

### 💰 Sales Management
- Sales transaction recording
- Sales analytics and reporting
- Product sales tracking
- Revenue calculations
- Sales data export

### 📚 API Documentation
- **Swagger/OpenAPI** documentation
- Interactive API testing
- Comprehensive endpoint documentation
- Request/response examples

## 🛠️ Tech Stack

- **Backend**: Go 1.21+ with Gin framework
- **Database**: MongoDB with official Go driver
- **Authentication**: JWT tokens
- **File Upload**: Multipart form-data handling
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Architecture**: Clean Architecture pattern

## 📁 Project Structure

```
agricultural-backend/
├── cmd/
│   ├── main.go              # Application entry point
│   └── seed/                # Database seeding
├── internal/
│   ├── config/              # Configuration management
│   ├── delivery/http/       # HTTP handlers & middleware
│   ├── domain/              # Domain models & interfaces
│   ├── infrastructure/      # Database & external services
│   ├── repository/          # Data access layer
│   ├── usecase/             # Business logic
│   └── utils/               # Utility functions
├── docs/                    # Swagger documentation
├── uploads/                 # Uploaded images storage
├── docker-compose.yml       # Docker services
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker & Docker Compose
- MongoDB (or use Docker)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/chatchomphu1000/agricultural-backend.git
   cd agricultural-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start MongoDB with Docker**
   ```bash
   docker-compose up -d
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

6. **Seed the database (optional)**
   ```bash
   go run cmd/seed/main.go
   ```

## 📖 API Documentation

Once the server is running, visit:
- **Swagger UI**: http://localhost:8082/swagger/index.html
- **Health Check**: http://localhost:8082/health

## 🔧 Configuration

Key environment variables:

```env
# Database
MONGODB_URI=mongodb://root:example@localhost:27017/Agricultural?authSource=admin
MONGODB_DATABASE=Agricultural

# JWT
JWT_SECRET=your-super-secret-jwt-key-here

# Server
PORT=8082

# Frontend (CORS)
FRONTEND_URL=http://localhost:3000

# Admin User (for seeding)
ADMIN_EMAIL=admin@agricultural.com
ADMIN_PASSWORD=password123
```

## 🖼️ Image Upload Features

### Supported Methods
1. **Local File Upload** - Multipart form-data
2. **Image URLs** - Direct URL links
3. **Mixed Upload** - Both files and URLs together

### Example Usage

#### Upload Files with Form Data
```bash
curl -X POST "http://localhost:8082/api/products" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "name=Tractor" \
  -F "price=25000" \
  -F "category=Tractors" \
  -F "stock=5" \
  -F "images=@image1.jpg" \
  -F "images=@image2.jpg" \
  -F "image_urls=https://example.com/image3.jpg"
```

#### JSON with Image URLs
```bash
curl -X POST "http://localhost:8082/api/products" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Plow",
    "price": 5000,
    "category": "Plows",
    "stock": 10,
    "image_urls": [
      "https://example.com/plow1.jpg",
      "https://example.com/plow2.jpg"
    ]
  }'
```

## 📊 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### Products
- `GET /api/products` - List products (public)
- `GET /api/products/{id}` - Get product details
- `POST /api/products` - Create product (admin)
- `PUT /api/products/{id}` - Update product (admin)
- `DELETE /api/products/{id}` - Delete product (admin)

### Inventory
- `PUT /api/inventories/{id}/stock` - Update stock (admin)
- `GET /api/inventories/low-stock` - Low stock products (admin)
- `GET /api/inventories/summary` - Inventory summary (admin)

### Sales
- `POST /api/sales` - Record sale (admin)
- `GET /api/sales` - List sales (admin)
- `GET /api/sales/summary` - Sales summary (admin)

### Categories
- `GET /api/categories` - List categories (public)
- `POST /api/categories` - Create category (admin)

## 🐳 Docker Support

Run with Docker Compose:
```bash
docker-compose up -d
```

This starts:
- MongoDB database
- MongoDB Express (admin interface)
- Application server

## 🧪 Testing

The project includes test utilities:
- `test_image_upload.html` - Interactive web interface
- Swagger UI for API testing
- Example curl commands

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

For support, please open an issue in the GitHub repository.

---

**Built with ❤️ for agricultural equipment management**
