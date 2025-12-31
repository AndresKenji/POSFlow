# Claude AI Instructions for POSFlow

## Project Overview
POSFlow is an offline-first Point of Sale system for restaurants and retail stores. It's built with Electron (frontend) and Go (backend), using SQLite for local data storage.

## Technology Stack
- **Frontend**: Electron, HTML, CSS, Vanilla JavaScript
- **Backend**: Go 1.21+, Gin/Fiber framework
- **Database**: SQLite
- **Architecture**: Offline-first, REST API, single executable distribution

## Project Goals
1. **Simplicity**: Keep code clean and maintainable
2. **Performance**: Fast responses, minimal latency
3. **Reliability**: Works 100% offline, data integrity is critical
4. **Usability**: Intuitive interfaces for non-technical users

## Code Style Guidelines

### Go Backend
- Follow standard Go conventions (gofmt, golint)
- Use descriptive variable names (e.g., `totalSales` not `ts`)
- Keep functions small and focused (single responsibility)
- Use proper error handling (never ignore errors)
- Add comments for exported functions and types
- Use struct tags for JSON serialization
- Organize code in packages: handlers, services, models, database

Example:
```go
// GetOrderByID retrieves an order by its ID
func GetOrderByID(orderID int) (*Order, error) {
    var order Order
    err := db.QueryRow(
        "SELECT * FROM orders WHERE id = ?",
        orderID,
    ).Scan(&order)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get order: %w", err)
    }

    return &order, nil
}
```

### JavaScript Frontend
- Use ES6+ features (const/let, arrow functions, async/await)
- Avoid jQuery or heavy frameworks (keep it lightweight)
- Use descriptive variable names in camelCase
- Add JSDoc comments for complex functions
- Handle errors gracefully with try-catch
- Always validate user input before sending to API

Example:
```javascript
/**
 * Creates a new order
 * @param {Object} orderData - The order information
 * @returns {Promise<Object>} The created order
 */
async function createOrder(orderData) {
    try {
        // Implementation here
    } catch (error) {
        console.error('Failed to create order:', error);
        throw error;
    }
}
```

### Database
- Use snake_case for table and column names
- Always use prepared statements (database/sql package handles this)
- Create indexes for frequently queried columns
- Include created_at and updated_at timestamps on all tables
- Use foreign keys to maintain referential integrity
- Use sqlx or GORM for easier database operations

## Architecture Patterns

### API Structure
- RESTful endpoints (GET, POST, PUT, DELETE)
- Consistent response format:
```json
{
    "success": true,
    "data": {...},
    "message": "Operation completed successfully"
}
```
- Use HTTP status codes correctly (200, 201, 400, 404, 500)
- Validate all inputs with struct validation (using validator package)

### Frontend-Backend Communication
- All API calls go through `api-client.js`
- Use async/await for all API calls
- Show loading states during operations
- Display user-friendly error messages
- Implement retry logic for failed requests

### State Management
- Keep state in memory (no localStorage needed)
- Use WebSockets for real-time updates between views
- Refresh data after mutations

## Key Features to Remember

### 1. Inventory Management
- Track stock levels in real-time
- Alert when stock is low
- Record all inventory movements (add, remove, sale)

### 2. Order Management
- Orders processed in FIFO (first-in, first-out) order
- Status flow: pending → preparing → ready → completed
- Kitchen view updates automatically when new orders arrive

### 3. Daily Menu
- Menu can be updated daily by admin
- Items can be marked as "out of stock" without deleting
- Menu items link to inventory products

### 4. Sales Reports
- End-of-day closing generates sales report
- Export to CSV/Excel for accounting
- Track daily, weekly, monthly trends

### 5. User Roles
- **Admin**: Full access, can modify everything
- **Kitchen**: View orders, update order status
- **Cashier**: Create orders, process payments

## Security Considerations
- Hash passwords with bcrypt (golang.org/x/crypto/bcrypt)
- Use JWT tokens for session management (github.com/golang-jwt/jwt)
- Sanitize all user inputs
- Validate file uploads (if implemented)
- Keep SQLite database file secure
- Use Go's built-in SQL injection protection via prepared statements

## Performance Guidelines
- Database queries should complete in <100ms
- UI should respond to user actions in <50ms
- Optimize images (use WebP format, compress)
- Lazy load data when possible
- Use database indexes for common queries

## Testing Requirements
- Write unit tests for all service functions
- Test API endpoints with different inputs
- Test error handling and edge cases
- Maintain >80% code coverage

## Common Patterns

### Adding a New Feature
1. Create database model struct (`backend/models/`)
2. Create request/response DTOs (`backend/dto/`)
3. Create service with business logic (`backend/services/`)
4. Create API handlers (`backend/handlers/`)
5. Register routes (`backend/routes/`)
6. Create frontend UI (`electron-app/src/views/`)
7. Add API client functions (`electron-app/src/js/api-client.js`)
8. Wire up frontend logic
9. Write tests

### Error Handling Pattern
```go
// Backend
func PerformOperation(c *gin.Context) {
    result, err := performOperation()
    if err != nil {
        log.Printf("Error performing operation: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": "Internal server error",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": result,
        "message": "Operation completed successfully",
    })
}
```

```javascript
// Frontend
try {
    const response = await apiClient.createOrder(orderData);
    showSuccessMessage('Order created successfully');
    updateUI(response.data);
} catch (error) {
    showErrorMessage('Failed to create order. Please try again.');
    console.error(error);
}
```

## When Helping with Code

### DO:
- ✅ Provide complete, working code examples
- ✅ Include error handling
- ✅ Add comments explaining complex logic
- ✅ Consider edge cases
- ✅ Follow the established patterns in the codebase
- ✅ Write code that works offline
- ✅ Optimize for performance

### DON'T:
- ❌ Use external APIs that require internet
- ❌ Add dependencies without justification
- ❌ Write overly complex solutions
- ❌ Ignore error handling
- ❌ Use deprecated libraries
- ❌ Add features that require cloud services

## Debugging Tips
- Check backend logs (use log package)
- Use Electron DevTools (Ctrl+Shift+I)
- Verify SQLite database with DB Browser
- Test API endpoints with curl or Postman
- Check CORS settings if API calls fail
- Use `go run -race` to detect race conditions
- Use pprof for performance profiling

## Build and Distribution
- Build backend: `go build -ldflags="-s -w" -o posflow-server`
- Use UPX to compress executable (optional): `upx --best posflow-server`
- Bundle backend executable with Electron app
- Use electron-builder for cross-platform distribution
- Backend binary should be placed in app resources folder

## Resources
- Go docs: https://go.dev/doc/
- Gin framework: https://gin-gonic.com/docs/
- Electron docs: https://www.electronjs.org/docs
- SQLite docs: https://www.sqlite.org/docs.html
- GORM (Go ORM): https://gorm.io/docs/

## Questions to Ask Before Implementing
1. Does this work offline?
2. Is this the simplest solution?
3. Have I handled errors?
4. Is this performant?
5. Will this scale to 1000+ products/orders?
6. Is this user-friendly for non-technical people?

---

**Remember**: POSFlow is designed for small businesses. Prioritize reliability, simplicity, and offline functionality over fancy features.