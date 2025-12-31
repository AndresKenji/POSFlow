# POSFlow - Domain-Driven Design Structure

## DDD Overview for POSFlow

Domain-Driven Design organizes code around **business domains** rather than technical layers. For POSFlow, we have these domains:

1. **Product** - Manages menu items and inventory
2. **Order** - Handles customer orders and their lifecycle
3. **Sales** - Manages daily sales and reporting
4. **Inventory** - Tracks stock movements

## Complete Project Structure

```
posflow/
├── electron-app/                    # Frontend (already created)
│
└── backend/                         # Go Backend with DDD
    ├── cmd/
    │   └── api/
    │       └── main.go             # Application entry point
    │
    ├── internal/
    │   ├── domain/                 # DOMAIN LAYER (Business Logic)
    │   │   ├── product/
    │   │   │   ├── entity.go       # Product entity
    │   │   │   ├── repository.go   # Product repository interface
    │   │   │   ├── service.go      # Product domain service
    │   │   │   └── value_objects.go # Money, Category, etc.
    │   │   │
    │   │   ├── order/
    │   │   │   ├── entity.go       # Order entity
    │   │   │   ├── repository.go   # Order repository interface
    │   │   │   ├── service.go      # Order domain service
    │   │   │   └── value_objects.go # OrderStatus, TableNumber
    │   │   │
    │   │   ├── inventory/
    │   │   │   ├── entity.go       # InventoryMovement entity
    │   │   │   ├── repository.go   # Inventory repository interface
    │   │   │   └── service.go      # Inventory domain service
    │   │   │
    │   │   ├── sales/
    │   │   │   ├── entity.go       # DailySales entity
    │   │   │   ├── repository.go   # Sales repository interface
    │   │   │   └── service.go      # Sales domain service
    │   │   │
    │   │   └── shared/             # Shared domain concepts
    │   │       ├── errors.go       # Domain errors
    │   │       └── events.go       # Domain events
    │   │
    │   ├── application/            # APPLICATION LAYER (Use Cases)
    │   │   ├── product/
    │   │   │   ├── commands/       # Write operations
    │   │   │   │   ├── create_product.go
    │   │   │   │   ├── update_product.go
    │   │   │   │   ├── delete_product.go
    │   │   │   │   └── update_stock.go
    │   │   │   ├── queries/        # Read operations
    │   │   │   │   ├── get_product.go
    │   │   │   │   ├── list_products.go
    │   │   │   │   └── get_low_stock.go
    │   │   │   └── dto/            # Data Transfer Objects
    │   │   │       └── product_dto.go
    │   │   │
    │   │   ├── order/
    │   │   │   ├── commands/
    │   │   │   │   ├── create_order.go
    │   │   │   │   ├── update_order_status.go
    │   │   │   │   └── cancel_order.go
    │   │   │   ├── queries/
    │   │   │   │   ├── get_order.go
    │   │   │   │   ├── list_orders.go
    │   │   │   │   └── get_pending_orders.go
    │   │   │   └── dto/
    │   │   │       └── order_dto.go
    │   │   │
    │   │   ├── sales/
    │   │   │   ├── commands/
    │   │   │   │   └── close_day.go
    │   │   │   └── queries/
    │   │   │       ├── get_daily_sales.go
    │   │   │       └── get_sales_report.go
    │   │   │
    │   │   └── inventory/
    │   │       ├── commands/
    │   │       │   └── record_movement.go
    │   │       └── queries/
    │   │           └── get_movements.go
    │   │
    │   ├── infrastructure/         # INFRASTRUCTURE LAYER
    │   │   ├── persistence/        # Database implementations
    │   │   │   ├── sqlite/
    │   │   │   │   ├── connection.go
    │   │   │   │   ├── product_repository.go
    │   │   │   │   ├── order_repository.go
    │   │   │   │   ├── inventory_repository.go
    │   │   │   │   ├── sales_repository.go
    │   │   │   │   └── models.go   # GORM models
    │   │   │   └── migrations/
    │   │   │       └── migrations.go
    │   │   │
    │   │   ├── http/               # HTTP server
    │   │   │   ├── server.go
    │   │   │   ├── middleware/
    │   │   │   │   ├── cors.go
    │   │   │   │   ├── logger.go
    │   │   │   │   └── error_handler.go
    │   │   │   └── handlers/
    │   │   │       ├── product_handler.go
    │   │   │       ├── order_handler.go
    │   │   │       ├── sales_handler.go
    │   │   │       └── health_handler.go
    │   │   │
    │   │   └── config/             # Configuration
    │   │       └── config.go
    │   │
    │   └── interfaces/             # INTERFACES LAYER (API)
    │       └── api/
    │           └── rest/
    │               └── router.go   # Routes definition
    │
    ├── pkg/                        # Shared packages
    │   ├── logger/
    │   │   └── logger.go
    │   └── validator/
    │       └── validator.go
    │
    ├── database/                   # Database files
    │   └── posflow.db
    │
    ├── go.mod
    ├── go.sum
    ├── Makefile
    └── README.md
```

## DDD Layers Explained

### 1. Domain Layer (Core Business Logic)
- **Entities**: Business objects with identity (Product, Order)
- **Value Objects**: Immutable objects (Money, OrderStatus)
- **Domain Services**: Business logic that doesn't belong to entities
- **Repository Interfaces**: Define how to persist data (implementation in infrastructure)

### 2. Application Layer (Use Cases)
- **Commands**: Operations that change state (CreateOrder, UpdateStock)
- **Queries**: Operations that read data (GetProduct, ListOrders)
- **DTOs**: Data structures for input/output

### 3. Infrastructure Layer (Technical Details)
- **Persistence**: Database implementations (SQLite with GORM)
- **HTTP**: Web server and handlers
- **Config**: Application configuration

### 4. Interfaces Layer (External Communication)
- **REST API**: HTTP endpoints
- **Routes**: API route definitions

## Key DDD Concepts for POSFlow

### Entities vs Value Objects

**Entity** (has identity):
```go
type Product struct {
    ID          ProductID  // Unique identifier
    Name        string
    Price       Money
    Stock       int
}
```

**Value Object** (no identity, immutable):
```go
type Money struct {
    Amount   float64
    Currency string
}

type OrderStatus string
const (
    StatusPending   OrderStatus = "pending"
    StatusPreparing OrderStatus = "preparing"
    StatusReady     OrderStatus = "ready"
    StatusCompleted OrderStatus = "completed"
)
```

### Aggregates

An **Aggregate** is a cluster of domain objects that can be treated as a single unit.

**Order Aggregate**:
```
Order (Aggregate Root)
  ├── OrderItems
  ├── TableNumber
  ├── Status
  └── Total
```

**Product Aggregate**:
```
Product (Aggregate Root)
  ├── Name
  ├── Price
  ├── Stock
  └── Category
```

### Repository Pattern

Repositories provide an abstraction for data access:

```go
// Domain layer defines the interface
type ProductRepository interface {
    Save(product *Product) error
    FindByID(id ProductID) (*Product, error)
    FindAll() ([]*Product, error)
    Delete(id ProductID) error
}

// Infrastructure layer implements it
type SQLiteProductRepository struct {
    db *gorm.DB
}
```

## Step-by-Step Implementation Guide

### Step 1: Setup Project Structure

```bash
cd posflow/backend

# Create directory structure
mkdir -p internal/{domain,application,infrastructure,interfaces}
mkdir -p internal/domain/{product,order,inventory,sales,shared}
mkdir -p internal/application/{product,order,sales,inventory}/{commands,queries,dto}
mkdir -p internal/infrastructure/{persistence/sqlite,http/{handlers,middleware},config}
mkdir -p internal/interfaces/api/rest
mkdir -p pkg/{logger,validator}
mkdir -p cmd/api
mkdir -p database
```

### Step 2: Initialize Go Module

```bash
go mod init github.com/yourusername/posflow

# Install dependencies
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/google/uuid
```

### Step 3: Create Domain Entities

We'll start with the **Product domain** as an example.

### Step 4: Create Repository Interfaces

Define how to persist data (interface only, implementation comes later).

### Step 5: Create Application Use Cases

Implement business use cases (commands and queries).

### Step 6: Create Infrastructure

Implement repositories with SQLite, setup HTTP server.

### Step 7: Wire Everything Together

Connect all layers in main.go using dependency injection.

## Benefits of DDD for POSFlow

✅ **Clear separation of concerns**
✅ **Business logic is independent of database/framework**
✅ **Easy to test** (mock repositories)
✅ **Scalable** (add new features without breaking existing ones)
✅ **Maintainable** (each domain is independent)
✅ **Database agnostic** (can switch from SQLite to PostgreSQL easily)

## Flow Example: Creating an Order

```
1. HTTP Request arrives
   ↓
2. Handler (infrastructure/http/handlers/order_handler.go)
   ↓
3. Command (application/order/commands/create_order.go)
   ↓
4. Domain Service (domain/order/service.go)
   - Validates business rules
   - Creates Order entity
   ↓
5. Repository (infrastructure/persistence/sqlite/order_repository.go)
   - Persists to database
   ↓
6. Response back through the layers
```

## Next Steps

Let me guide you through implementing each layer step by step:

1. **Domain Layer** - Start with Product domain
2. **Application Layer** - Create use cases
3. **Infrastructure Layer** - Implement persistence
4. **Interfaces Layer** - Setup HTTP handlers
5. **Main** - Wire everything together

Ready to start? Let's begin with the Domain Layer!
