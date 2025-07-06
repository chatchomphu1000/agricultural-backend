<<<<<<< HEAD
# Agricultural Equipment Store API

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![MongoDB](https://img.shields.io/badge/MongoDB-6.0+-green.svg)](https://mongodb.com)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A comprehensive REST API for managing agricultural equipment store operations, built with Go using Clean Architecture principles.

## 🚀 Features

### Core Features
- **Product Management** - Complete CRUD operations for agricultural equipment
- **Multiple Image Support** - Upload files or provide URLs, support multiple images per product
- **Inventory Management** - Track stock levels, low stock alerts, and inventory summaries
- **Sales Management** - Record sales transactions and generate reports
- **Category Management** - Organize products by categories
- **User Authentication** - JWT-based authentication with admin and user roles

### API Capabilities
- **RESTful API** with comprehensive endpoints
- **Swagger Documentation** - Auto-generated API documentation
- **File Upload Support** - Handle product images with validation
- **Filtering & Pagination** - Advanced search and pagination features
- **CORS Support** - Cross-origin resource sharing enabled
- **Error Handling** - Comprehensive error responses

### Technical Features
- **Clean Architecture** - Separation of concerns with domain-driven design
- **MongoDB Integration** - NoSQL database for flexible data storage
- **Docker Support** - Containerized deployment ready
- **Environment Configuration** - Flexible configuration management
- **Middleware Support** - Authentication, logging, and CORS middleware
- **Static File Serving** - Serve uploaded images directly

## Tech Stack

- **Backend**: Go 1.21, Gin Web Framework
- **Database**: MongoDB 7.0
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Database Admin**: Mongo Express

## Project Structure

```
backend-new/
├── cmd/
│   ├── main.go              # Application entry point
│   └── seed/
│       └── main.go          # Database seeder
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── domain/
│   │   ├── user.go          # User domain models
│   │   ├── product.go       # Product domain models
│   │   └── repository.go    # Repository interfaces
│   ├── usecase/
│   │   ├── auth_usecase.go  # Authentication business logic
│   │   └── product_usecase.go # Product business logic
│   ├── repository/
│   │   ├── user_repository.go    # User data operations
│   │   └── product_repository.go # Product data operations
│   ├── delivery/
│   │   └── http/
│   │       ├── server.go         # HTTP server setup
│   │       ├── auth_handler.go   # Authentication handlers
│   │       ├── product_handler.go # Product handlers
│   │       └── middleware/
│   │           └── auth_middleware.go # JWT middleware
│   └── infrastructure/
│       ├── database/
│       │   └── mongodb.go   # MongoDB connection
│       └── logger/
│           └── logger.go    # Logging utility
├── docs/
│   └── docs.go              # Swagger documentation
├── docker-compose.yml       # Docker services configuration
├── init-mongo.js           # MongoDB initialization script
├── .env                    # Environment variables
├── go.mod                  # Go module dependencies
└── README.md              # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for database)
- Git

### Installation

1. **Clone the repository**
   ```powershell
   git clone <repository-url>
   cd backend-new
   ```

2. **Start MongoDB with Docker**
   ```powershell
   docker-compose up -d
   ```

3. **Install Go dependencies**
   ```powershell
   go mod download
   ```

4. **Set up environment variables**
   - Copy `.env.example` to `.env`
   - Update the configuration as needed

5. **Run database migrations and seed data**
   ```powershell
   go run cmd/seed/main.go
   ```

6. **Start the application**
   ```powershell
   go run cmd/main.go
   ```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/profile` - Get user profile (requires authentication)

### Products
- `GET /api/products` - Get all products (public)
- `GET /api/products/:id` - Get product by ID (public)
- `POST /api/products` - Create product (admin only)
- `PUT /api/products/:id` - Update product (admin only)
- `DELETE /api/products/:id` - Delete product (admin only)

### Other
- `GET /health` - Health check
- `GET /swagger/index.html` - API documentation

## Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Default Admin User
- Email: `admin@agricultural.com`
- Password: `password123`

## Database

The application uses MongoDB with the following collections:
- `users` - User accounts and authentication
- `products` - Agricultural equipment products

### Database Management

Access MongoDB through:
- **Mongo Express**: http://localhost:8081
- **Direct MongoDB**: mongodb://root:example@localhost:27017

## Environment Variables

```env
# Database Configuration
MONGODB_URI=mongodb://root:example@localhost:27017
MONGODB_DATABASE=Agricultural

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# Server Configuration
PORT=8082

# Frontend URL
FRONTEND_URL=http://localhost:3000

# Admin User
ADMIN_EMAIL=admin@agricultural.com
ADMIN_PASSWORD=password123
```

## API Documentation

Once the server is running, visit:
- Swagger UI: http://localhost:8082/swagger/index.html

## Development

### Running Tests
```powershell
go test ./...
```

### Building for Production
```powershell
go build -o bin/server cmd/main.go
```

### Docker Build
```powershell
docker build -t agricultural-api .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.

## Support

For support, please contact the development team or create an issue in the repository.
=======
# agricultural-backend
>>>>>>> cdad1bb4f8d510b0c945fed6a87547815392a75b
