# Backend Setup Guide (Go)

This guide will help you set up the Go backend for POSFlow.

## Prerequisites

- Go 1.21 or higher
- SQLite3
- Basic knowledge of Go and REST APIs

## Quick Start

### 1. Initialize Go Module

```bash
cd backend
go mod init github.com/yourusername/posflow
```

### 2. Install Dependencies

```bash
# Web framework (choose one)
go get -u github.com/gin-gonic/gin
# OR
go get -u github.com/gofiber/fiber/v2

# Database
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite

# Security
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/golang-jwt/jwt/v5

# Validation
go get -u github.com/go-playground/validator/v10

# Environment variables
go get -u github.com/joho/godotenv

# CORS
go get -u github.com/gin-contrib/cors
# OR for Fiber
go get -u github.com/gofiber/fiber/v2/middleware/cors
```

### 3. Create Project Structure

```bash
mkdir -p backend/cmd/server
mkdir -p backend/internal/handlers
mkdir -p backend/internal/services
mkdir -p backend/internal/models
mkdir -p backend/internal/database
mkdir -p backend/internal/middleware
mkdir -p backend/internal/dto
mkdir -p backend/config
mkdir -p backend/logs
```

### 4. Environment Configuration

Create a `.env` file in the `backend/` directory:

```env
# Server Configuration
PORT=8080
ENV=development

# Database
DB_PATH=./posflow.db

# JWT Secret
JWT_SECRET=your-secret-key-change-this-in-production

# CORS
CORS_ORIGINS=*

# Logging
LOG_LEVEL=info
LOG_FILE=./logs/posflow.log
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── handlers/                   # HTTP handlers (controllers)
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   ├── inventory_handler.go
│   │   ├── sales_handler.go
│   │   ├── auth_handler.go
│   │   └── menu_handler.go
│   ├── services/                   # Business logic
│   │   ├── order_service.go
│   │   ├── product_service.go
│   │   ├── inventory_service.go
│   │   ├── sales_service.go
│   │   └── auth_service.go
│   ├── models/                     # Database models
│   │   ├── order.go
│   │   ├── product.go
│   │   ├── user.go
│   │   └── inventory.go
│   ├── dto/                        # Data Transfer Objects
│   │   ├── order_dto.go
│   │   ├── product_dto.go
│   │   └── auth_dto.go
│   ├── database/                   # Database connection and migrations
│   │   ├── db.go
│   │   └── migrations.go
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── utils/                      # Utility functions
│       ├── response.go
│       └── validator.go
├── config/
│   └── config.go                   # Configuration management
├── logs/
│   └── .gitkeep
├── .env
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## Implementation Steps

### Step 1: Basic Server Setup

Create `cmd/server/main.go`:

```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/posflow/internal/database"
    "github.com/yourusername/posflow/config"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Initialize database
    db, err := database.InitDB(cfg.DBPath)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Run migrations
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Setup router
    router := gin.Default()

    // CORS middleware
    router.Use(corsMiddleware())

    // Health check
    router.GET("/api/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "success": true,
            "data": gin.H{
                "status": "ok",
                "version": "1.0.0",
            },
        })
    })

    // Setup routes
    setupRoutes(router, db)

    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### Step 2: Database Setup

Create `internal/database/db.go`:

```go
package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

Create `internal/database/migrations.go`:

```go
package database

import (
    "gorm.io/gorm"
    "github.com/yourusername/posflow/internal/models"
)

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Product{},
        &models.Order{},
        &models.OrderItem{},
        &models.InventoryMovement{},
        &models.Sale{},
    )
}
```

### Step 3: Define Models

Create `internal/models/product.go`:

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"not null" json:"name"`
    Description string         `json:"description"`
    Price       float64        `gorm:"not null" json:"price"`
    Category    string         `json:"category"`
    Available   bool           `gorm:"default:true" json:"available"`
    Stock       int            `gorm:"default:0" json:"stock"`
    ImageURL    string         `json:"imageUrl"`
    CreatedAt   time.Time      `json:"createdAt"`
    UpdatedAt   time.Time      `json:"updatedAt"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Step 4: Create Response Utility

Create `internal/utils/response.go`:

```go
package utils

import "github.com/gin-gonic/gin"

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}, message string) {
    c.JSON(code, Response{
        Success: true,
        Data:    data,
        Message: message,
    })
}

func ErrorResponse(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Success: false,
        Message: message,
    })
}
```

### Step 5: Create Handlers

Create `internal/handlers/product_handler.go`:

```go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/yourusername/posflow/internal/services"
    "github.com/yourusername/posflow/internal/utils"
)

type ProductHandler struct {
    service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
    products, err := h.service.GetAll()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessResponse(c, http.StatusOK, gin.H{
        "products": products,
    }, "")
}

func (h *ProductHandler) GetByID(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

    product, err := h.service.GetByID(uint(id))
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
        return
    }

    utils.SuccessResponse(c, http.StatusOK, product, "")
}

// ... more methods
```

### Step 6: Setup Routes

Create route setup in `main.go`:

```go
func setupRoutes(router *gin.Engine, db *gorm.DB) {
    api := router.Group("/api")

    // Initialize services
    productService := services.NewProductService(db)
    orderService := services.NewOrderService(db)
    // ... more services

    // Initialize handlers
    productHandler := handlers.NewProductHandler(productService)
    orderHandler := handlers.NewOrderHandler(orderService)
    // ... more handlers

    // Product routes
    products := api.Group("/products")
    {
        products.GET("", productHandler.GetAll)
        products.GET("/:id", productHandler.GetByID)
        products.POST("", productHandler.Create)
        products.PUT("/:id", productHandler.Update)
        products.DELETE("/:id", productHandler.Delete)
    }

    // Order routes
    orders := api.Group("/orders")
    {
        orders.GET("", orderHandler.GetAll)
        orders.GET("/:id", orderHandler.GetByID)
        orders.POST("", orderHandler.Create)
        orders.PUT("/:id/status", orderHandler.UpdateStatus)
        orders.DELETE("/:id", orderHandler.Delete)
    }

    // ... more routes
}
```

## Building the Application

### Development

```bash
cd backend
go run cmd/server/main.go
```

### Production Build

```bash
# Build with optimizations
go build -ldflags="-s -w" -o posflow-server cmd/server/main.go

# Optional: Compress with UPX
upx --best posflow-server
```

### Cross-Platform Build

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o posflow-server.exe cmd/server/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o posflow-server-mac cmd/server/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o posflow-server-linux cmd/server/main.go
```

## Testing

### Run Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Deployment

### 1. Build the Binary

```bash
go build -ldflags="-s -w" -o posflow-server cmd/server/main.go
```

### 2. Package with Electron

The built binary should be placed in `electron-app/resources/` so Electron can spawn it.

### 3. Electron Integration

In `electron-app/main.js`, add code to spawn the Go backend:

```javascript
const { spawn } = require('child_process');
const path = require('path');

let backendProcess;

function startBackend() {
  const backendPath = path.join(
    process.resourcesPath,
    'posflow-server' + (process.platform === 'win32' ? '.exe' : '')
  );

  backendProcess = spawn(backendPath);

  backendProcess.stdout.on('data', (data) => {
    console.log(`Backend: ${data}`);
  });

  backendProcess.stderr.on('data', (data) => {
    console.error(`Backend Error: ${data}`);
  });
}

app.on('ready', () => {
  startBackend();
  setTimeout(createMainWindow, 2000); // Wait for backend to start
});

app.on('quit', () => {
  if (backendProcess) {
    backendProcess.kill();
  }
});
```

## Troubleshooting

### Database Locked Error

If you get a "database is locked" error:
- Ensure only one instance of the server is running
- Check file permissions on the database file
- Use `PRAGMA busy_timeout = 5000` in database connection

### CORS Issues

Make sure CORS middleware is properly configured to allow requests from the Electron app.

### Port Already in Use

If port 8080 is already in use, change it in the `.env` file.

## Next Steps

1. Implement all handlers according to [API_SPEC.md](../API_SPEC.md)
2. Add authentication middleware
3. Implement JWT token generation and validation
4. Add input validation using validator package
5. Implement comprehensive error handling
6. Add logging throughout the application
7. Write unit tests for services
8. Write integration tests for handlers
9. Add database seeding for initial data
10. Optimize database queries with indexes

## Useful Commands

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Update dependencies
go get -u ./...
go mod tidy

# View module dependencies
go mod graph

# Check for security vulnerabilities
go list -json -m all | nancy sleuth
```

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [API Specification](../API_SPEC.md)
- [Project Documentation](../CLAUDE.md)
