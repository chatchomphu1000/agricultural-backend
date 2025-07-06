# Claude Prompts for Agricultural Equipment Store API

## Project Overview
This is a Go-based REST API for an Agricultural Equipment Store built with Clean Architecture principles.

## Technology Stack
- **Language**: Go 1.21
- **Framework**: Gin (HTTP router)
- **Database**: MongoDB
- **Authentication**: JWT
- **Documentation**: Swagger/OpenAPI
- **Architecture**: Clean Architecture

## Project Structure
```
agricultural-equipment-store/
├── cmd/
│   ├── main.go              # Application entry point
│   └── seed/
│       └── main.go          # Database seeding
├── docs/                    # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── config/              # Configuration management
│   │   └── config.go
│   ├── delivery/            # HTTP handlers (controllers)
│   │   └── http/
│   │       ├── auth_handler.go
│   │       ├── product_handler.go
│   │       ├── server.go
│   │       └── middleware/
│   │           └── auth_middleware.go
│   ├── domain/              # Business entities
│   │   ├── product.go
│   │   ├── repository.go
│   │   └── user.go
│   ├── infrastructure/      # External dependencies
│   │   ├── database/
│   │   │   └── mongodb.go
│   │   └── logger/
│   │       └── logger.go
│   ├── repository/          # Data access layer
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   └── usecase/             # Business logic
│       ├── auth_usecase.go
│       └── product_usecase.go
├── docker-compose.yml       # Docker services
├── go.mod                   # Go modules
├── go.sum                   # Go dependencies
├── init-mongo.js           # MongoDB initialization
└── README.md               # Project documentation
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/profile` - Get user profile (protected)

### Products
- `GET /api/products` - Get all products (public)
- `GET /api/products/:id` - Get product by ID (public)
- `POST /api/products` - Create product (admin only)
- `PUT /api/products/:id` - Update product (admin only)
- `DELETE /api/products/:id` - Delete product (admin only)

### Inventory / Stock Management
- `PUT /api/products/:id/stock` - Update product stock (admin only)
- `GET /api/products/low-stock` - Get low stock products (admin only)
- `GET /api/stock/summary` - Get stock summary (admin only)

### Sales & Reports
- `POST /api/sales` - Create new sale (admin only)
- `GET /api/sales` - Get sales with filtering (admin only)
- `GET /api/sales/summary` - Get sales summary (admin only)
- `GET /api/sales/by-product` - Get sales by product (admin only)
- `GET /api/sales/export` - Export sales as CSV (admin only)

### Documentation
- `GET /swagger/index.html` - Swagger UI
- `GET /health` - Health check

## Common Prompts

### 1. Development Setup
```
Set up the development environment for this Go project:
1. Install dependencies: go mod download
2. Start MongoDB: docker-compose up -d
3. Run the application: go run ./cmd/main.go
4. Generate Swagger docs: swag init -g cmd/main.go -o docs
```

### 2. Database Operations
```
Help me with MongoDB operations:
1. Connect to MongoDB using the connection string in config
2. Create indexes for better performance
3. Implement CRUD operations for [entity]
4. Add data validation and error handling
```

### 3. API Development
```
Create a new API endpoint for [feature]:
1. Add domain models in internal/domain/
2. Create repository interface and implementation
3. Implement use case with business logic
4. Add HTTP handler with proper validation
5. Include Swagger documentation
6. Add authentication/authorization if needed
```

### 4. Authentication & Authorization
```
Implement JWT authentication:
1. Create JWT tokens on login
2. Validate tokens in middleware
3. Implement role-based access control
4. Add password hashing with bcrypt
5. Handle token expiration and refresh
```

### 5. Error Handling
```
Improve error handling:
1. Create custom error types
2. Add proper HTTP status codes
3. Implement error middleware
4. Add logging for debugging
5. Return consistent error responses
```

### 6. Testing
```
Add comprehensive testing:
1. Unit tests for use cases
2. Integration tests for repositories
3. HTTP handler tests
4. Mock dependencies for testing
5. Test coverage reports
```

### 7. Documentation
```
Update project documentation:
1. Generate/update Swagger documentation
2. Add API examples and usage
3. Update README with setup instructions
4. Document environment variables
5. Add deployment instructions
```

### 8. Performance Optimization
```
Optimize application performance:
1. Add database indexes
2. Implement caching strategy
3. Add request rate limiting
4. Optimize database queries
5. Add monitoring and metrics
```

### 9. Security Enhancements
```
Improve application security:
1. Add input validation and sanitization
2. Implement CORS properly
3. Add request size limits
4. Secure JWT implementation
5. Add security headers
```

### 10. Deployment
```
Prepare for deployment:
1. Create Docker configuration
2. Set up environment variables
3. Configure database migrations
4. Add health check endpoints
5. Set up monitoring and logging
```

## Environment Variables
```
# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=agricultural_store

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h

# Server
SERVER_PORT=8082
GIN_MODE=release

# Frontend
FRONTEND_URL=http://localhost:3000
```

## Common Commands
```bash
# Install dependencies
go mod download

# Run application
go run ./cmd/main.go

# Generate Swagger docs
swag init -g cmd/main.go -o docs

# Run tests
go test ./...

# Build application
go build -o bin/server ./cmd/main.go

# Start MongoDB
docker-compose up -d mongodb

# Seed database
go run ./cmd/seed/main.go
```

## Database Schema

### Users Collection
```json
{
  "_id": "ObjectId",
  "email": "string",
  "password": "string (hashed)",
  "name": "string",
  "role": "string (user/admin)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Products Collection
```json
{
  "_id": "ObjectId",
  "name": "string",
  "description": "string",
  "price": "number",
  "category": "string",
  "image_url": "string",
  "stock": "number",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Troubleshooting

### Common Issues
1. **Swagger not working**: Make sure docs package is imported in main.go
2. **Database connection failed**: Check MongoDB is running and connection string
3. **JWT token invalid**: Verify JWT secret and token format
4. **CORS errors**: Check CORS configuration in server.go
5. **Port already in use**: Change port in config or kill existing process
6. **Login returns 401**: Make sure admin user exists in database (run seed command)

### Authentication Issues
- **401 Unauthorized on login**: This usually means:
  - No users exist in the database (run `go run .\cmd\seed\main.go`)
  - Wrong email/password combination
  - Database connection issues
  - Password hashing/verification problems

### Test Login Credentials
- **Admin Email**: admin@agricultural.com
- **Admin Password**: password123
- **Test Command**: `Invoke-RestMethod -Uri "http://localhost:8082/api/auth/login" -Method Post -ContentType "application/json" -Body '{"email":"admin@agricultural.com","password":"password123"}'`

### Debug Commands
```bash
# Check if MongoDB is running
docker ps

# Check application logs
go run ./cmd/main.go 2>&1 | tee app.log

# Test API endpoints
curl -X GET http://localhost:8082/health
curl -X GET http://localhost:8082/swagger/index.html
```

## Future Enhancements
1. Add product categories and filtering
2. Implement shopping cart functionality
3. Add order management system
4. Implement payment integration
5. Add product reviews and ratings
6. Implement inventory management
7. Add email notifications
8. Implement file upload for product images
9. Add search functionality
10. Implement audit logging

## Code Quality Guidelines
1. Follow Go naming conventions
2. Use meaningful variable and function names
3. Add comments for complex logic
4. Implement proper error handling
5. Write unit tests for critical functions
6. Use dependency injection
7. Follow SOLID principles
8. Keep functions small and focused
9. Use interfaces for better testability
10. Validate input data properly

---

*Last updated: July 6, 2025*
