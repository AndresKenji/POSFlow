# POSFlow API Specification

This document defines the REST API that the Go backend must implement for POSFlow.

## Base Configuration

- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`
- **Response Format**: All responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully"
}
```

## Status Codes

- `200 OK`: Successful GET, PUT, DELETE
- `201 Created`: Successful POST
- `400 Bad Request`: Invalid input
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Endpoints

### Health Check

#### `GET /api/health`
Check if the server is running.

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "version": "1.0.0"
  }
}
```

---

## Orders

### `GET /api/orders`
Get all orders, optionally filtered by status.

**Query Parameters:**
- `status` (optional): Filter by status (`pending`, `preparing`, `ready`, `completed`)

**Response:**
```json
{
  "success": true,
  "data": {
    "orders": [
      {
        "id": 1,
        "orderNumber": "ORD-001",
        "items": [
          {
            "productId": 1,
            "name": "Burger",
            "quantity": 2,
            "price": 9.99
          }
        ],
        "total": 19.98,
        "status": "pending",
        "createdAt": "2025-01-01T10:30:00Z",
        "updatedAt": "2025-01-01T10:30:00Z"
      }
    ]
  }
}
```

### `GET /api/orders/:id`
Get a specific order by ID.

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "orderNumber": "ORD-001",
    "items": [...],
    "total": 19.98,
    "status": "pending",
    "createdAt": "2025-01-01T10:30:00Z",
    "updatedAt": "2025-01-01T10:30:00Z"
  }
}
```

### `POST /api/orders`
Create a new order.

**Request Body:**
```json
{
  "items": [
    {
      "productId": 1,
      "quantity": 2
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "orderNumber": "ORD-001",
    "items": [...],
    "total": 19.98,
    "status": "pending",
    "createdAt": "2025-01-01T10:30:00Z"
  },
  "message": "Order created successfully"
}
```

### `PUT /api/orders/:id/status`
Update order status.

**Request Body:**
```json
{
  "status": "preparing"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "status": "preparing",
    "updatedAt": "2025-01-01T10:35:00Z"
  },
  "message": "Order status updated"
}
```

### `DELETE /api/orders/:id`
Delete an order.

**Response:**
```json
{
  "success": true,
  "message": "Order deleted successfully"
}
```

---

## Products

### `GET /api/products`
Get all products.

**Query Parameters:**
- `available` (optional): Filter by availability (`true` or `false`)

**Response:**
```json
{
  "success": true,
  "data": {
    "products": [
      {
        "id": 1,
        "name": "Burger",
        "description": "Delicious beef burger",
        "price": 9.99,
        "category": "Main Course",
        "available": true,
        "stock": 50,
        "imageUrl": "/images/burger.jpg",
        "createdAt": "2025-01-01T00:00:00Z",
        "updatedAt": "2025-01-01T00:00:00Z"
      }
    ]
  }
}
```

### `GET /api/products/:id`
Get a specific product by ID.

### `POST /api/products`
Create a new product.

**Request Body:**
```json
{
  "name": "Burger",
  "description": "Delicious beef burger",
  "price": 9.99,
  "category": "Main Course",
  "available": true,
  "stock": 50,
  "imageUrl": "/images/burger.jpg"
}
```

### `PUT /api/products/:id`
Update a product.

**Request Body:** Same as POST

### `DELETE /api/products/:id`
Delete a product.

---

## Inventory

### `GET /api/inventory`
Get current inventory status for all products.

**Response:**
```json
{
  "success": true,
  "data": {
    "inventory": [
      {
        "productId": 1,
        "productName": "Burger",
        "currentStock": 50,
        "lowStockThreshold": 10,
        "isLowStock": false
      }
    ]
  }
}
```

### `POST /api/inventory`
Update inventory (add or remove stock).

**Request Body:**
```json
{
  "productId": 1,
  "quantity": -5,
  "reason": "Sale"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "productId": 1,
    "previousStock": 50,
    "newStock": 45,
    "change": -5
  },
  "message": "Inventory updated successfully"
}
```

---

## Sales

### `GET /api/sales/daily`
Get daily sales report.

**Query Parameters:**
- `date` (optional): Date in `YYYY-MM-DD` format (defaults to today)

**Response:**
```json
{
  "success": true,
  "data": {
    "date": "2025-01-01",
    "totalSales": 1250.50,
    "totalOrders": 45,
    "averageOrderValue": 27.79,
    "topProducts": [
      {
        "productId": 1,
        "name": "Burger",
        "quantitySold": 30,
        "revenue": 299.70
      }
    ]
  }
}
```

### `GET /api/sales/report`
Get sales report for a date range.

**Query Parameters:**
- `start`: Start date in `YYYY-MM-DD` format
- `end`: End date in `YYYY-MM-DD` format

**Response:**
```json
{
  "success": true,
  "data": {
    "startDate": "2025-01-01",
    "endDate": "2025-01-07",
    "totalSales": 8750.50,
    "totalOrders": 315,
    "averageOrderValue": 27.78,
    "dailyBreakdown": [
      {
        "date": "2025-01-01",
        "sales": 1250.50,
        "orders": 45
      }
    ]
  }
}
```

### `POST /api/sales/close-day`
Close the day and generate end-of-day report.

**Response:**
```json
{
  "success": true,
  "data": {
    "date": "2025-01-01",
    "totalSales": 1250.50,
    "totalOrders": 45,
    "cashAmount": 500.00,
    "cardAmount": 750.50,
    "closedAt": "2025-01-01T22:00:00Z"
  },
  "message": "Day closed successfully"
}
```

---

## Authentication

### `POST /api/auth/login`
User login.

**Request Body:**
```json
{
  "username": "admin",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin",
      "name": "Administrator"
    }
  },
  "message": "Login successful"
}
```

### `POST /api/auth/logout`
User logout.

**Response:**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

### `GET /api/auth/me`
Get current user information.

**Headers:**
- `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "role": "admin",
    "name": "Administrator"
  }
}
```

---

## Menu

### `GET /api/menu`
Get today's menu (available products).

**Response:**
```json
{
  "success": true,
  "data": {
    "date": "2025-01-01",
    "items": [
      {
        "id": 1,
        "name": "Burger",
        "description": "Delicious beef burger",
        "price": 9.99,
        "category": "Main Course",
        "available": true,
        "imageUrl": "/images/burger.jpg"
      }
    ]
  }
}
```

### `PUT /api/menu/:id/availability`
Update menu item availability.

**Request Body:**
```json
{
  "available": false
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "available": false
  },
  "message": "Menu item availability updated"
}
```

---

## Data Models

### Order
```go
type Order struct {
    ID          int64     `json:"id"`
    OrderNumber string    `json:"orderNumber"`
    Items       []OrderItem `json:"items"`
    Total       float64   `json:"total"`
    Status      string    `json:"status"` // pending, preparing, ready, completed
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type OrderItem struct {
    ProductID int64   `json:"productId"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}
```

### Product
```go
type Product struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Category    string    `json:"category"`
    Available   bool      `json:"available"`
    Stock       int       `json:"stock"`
    ImageURL    string    `json:"imageUrl"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

### User
```go
type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"` // Never send in response
    Role     string `json:"role"` // admin, kitchen, cashier
    Name     string `json:"name"`
}
```

---

## Error Response Format

When an error occurs, the response should follow this format:

```json
{
  "success": false,
  "message": "Error description here"
}
```

---

## CORS Configuration

The Go backend should enable CORS with the following settings:

```go
// Allow Electron app to access the API
AllowOrigins: ["*"] // Since it's a local app
AllowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
AllowHeaders: ["Content-Type", "Authorization"]
```

---

## Database Schema (SQLite)

```sql
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL,
    category TEXT,
    available BOOLEAN DEFAULT 1,
    stock INTEGER DEFAULT 0,
    image_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_number TEXT UNIQUE NOT NULL,
    total REAL NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    price REAL NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Inventory movements table
CREATE TABLE inventory_movements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Sales table
CREATE TABLE sales (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    payment_method TEXT,
    sale_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- Indexes
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_products_available ON products(available);
CREATE INDEX idx_sales_date ON sales(sale_date);
```

---

## Notes for Go Implementation

1. **Framework**: Use Gin or Fiber for the HTTP server
2. **Database**: Use `database/sql` with SQLite driver or GORM
3. **JWT**: Use `github.com/golang-jwt/jwt` for authentication
4. **Password Hashing**: Use `golang.org/x/crypto/bcrypt`
5. **Validation**: Use `github.com/go-playground/validator`
6. **Logging**: Use standard `log` package or `logrus`
7. **Configuration**: Use environment variables or config file

### Recommended Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   ├── auth_handler.go
│   │   └── ...
│   ├── services/
│   │   ├── order_service.go
│   │   ├── product_service.go
│   │   └── ...
│   ├── models/
│   │   ├── order.go
│   │   ├── product.go
│   │   └── ...
│   ├── database/
│   │   ├── db.go
│   │   └── migrations.go
│   └── middleware/
│       ├── auth.go
│       └── cors.go
├── go.mod
└── go.sum
```
