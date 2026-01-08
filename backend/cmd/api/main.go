package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// Application layer
	orderCommands "POSFlowBackend/internal/application/order/commands"
	orderQueries "POSFlowBackend/internal/application/order/querys"
	productCommands "POSFlowBackend/internal/application/product/commands"
	productQueries "POSFlowBackend/internal/application/product/queries"
	salesCommands "POSFlowBackend/internal/application/sales/commands"
	salesQueries "POSFlowBackend/internal/application/sales/queries"

	// Domain layer
	"POSFlowBackend/internal/domain/order"

	// Infrastructure layer
	"POSFlowBackend/internal/domain/sales"
	"POSFlowBackend/internal/infrastructure/config"
	"POSFlowBackend/internal/infrastructure/http"
	"POSFlowBackend/internal/infrastructure/http/handlers"
	"POSFlowBackend/internal/infrastructure/http/routes"
	"POSFlowBackend/internal/infrastructure/persistence/sqlite"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("📝 Configuration loaded - DB: %s, Port: %s", cfg.DatabasePath, cfg.ServerPort)

	// Initialize database
	database, err := sqlite.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Initialize repositories (Infrastructure layer)
	productRepo := sqlite.NewProductRepository(database.DB)
	orderRepo := sqlite.NewOrderRepository(database.DB)
	salesRepo := sqlite.NewSalesRepository(database.DB)
	log.Println("✅ Repositories initialized")

	// Initialize domain services
	orderService := order.NewOrderService(orderRepo, productRepo)
	salesService := sales.NewSalesService(salesRepo, orderRepo)
	log.Println("✅ Domain services initialized")

	// Initialize application layer - Product commands
	createProductCmd := productCommands.NewCreateProductCommand(productRepo)
	updateProductCmd := productCommands.NewUpdateProductCommand(productRepo)
	deleteProductCmd := productCommands.NewDeleteProductCommand(productRepo)
	updateStockCmd := productCommands.NewUpdateStockCommand(productRepo)

	// Initialize application layer - Product queries
	listProductsQuery := productQueries.NewListProductsQuery(productRepo)
	getProductQuery := productQueries.NewGetProductQuery(productRepo)
	getLowStockQuery := productQueries.NewGetLowStockQuery(productRepo)

	// Initialize application layer - Order commands
	createOrderCmd := orderCommands.NewCreateOrderCommand(orderService, productRepo)
	updateOrderStatusCmd := orderCommands.NewUpdateOrderStatusCommand(orderRepo, productRepo)

	// Initialize application layer - Order queries
	listOrdersQuery := orderQueries.NewListOrdersQuery(orderRepo, productRepo)
	getOrderQuery := orderQueries.NewGetOrderQuery(orderRepo, productRepo)
	getPendingOrdersQuery := orderQueries.NewGetPendingOrdersQuery(orderRepo, productRepo)

	// Initialize application layer - Sales commands
	closeDayCmd := salesCommands.NewCloseDayCommand(salesService)

	// Initialize application layer - Sales queries
	getDailySalesQuery := salesQueries.NewGetDailySalesQuery(salesService)
	getSalesReportQuery := salesQueries.NewGetSalesReportQuery(salesService)
	log.Println("✅ Application layer initialized")

	// Initialize HTTP handlers (Interfaces layer)
	productHandler := handlers.NewProductHandler(
		createProductCmd,
		updateProductCmd,
		deleteProductCmd,
		updateStockCmd,
		listProductsQuery,
		getProductQuery,
		getLowStockQuery,
	)

	orderHandler := handlers.NewOrderHandler(
		createOrderCmd,
		updateOrderStatusCmd,
		listOrdersQuery,
		getOrderQuery,
		getPendingOrdersQuery,
	)

	salesHandler := handlers.NewSalesHandler( // ← Agregar
		getDailySalesQuery,
		getSalesReportQuery,
		closeDayCmd,
	)

	log.Println("✅ HTTP handlers initialized")

	// Initialize HTTP server
	server := http.NewServer(cfg.ServerPort)

	// Register routes
	routes.RegisterRoutes(server.Router(), productHandler, orderHandler, salesHandler)
	log.Println("✅ Routes registered")

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("🛑 Shutting down server...")
	log.Println("✅ Server stopped gracefully")
}
