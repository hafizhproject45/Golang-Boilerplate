package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/config"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/database"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/middleware"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/route"
	"github.com/hafizhproject45/Golang-Boilerplate.git/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := setupFiberApp()
	db := setupDatabase()
	defer closeDatabase(db)
	rdb := setupRedis()
	defer rdb.Close()
	setupRoutes(app, db, rdb)

	address := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)

	// Start server and handle graceful shutdown
	serverErrors := make(chan error, 1)
	go startServer(app, address, serverErrors)
	handleGracefulShutdown(ctx, app, serverErrors)
}

func setupRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://redis:6379/0"
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		utils.Log.Fatalf("Redis URL parse error: %v", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		utils.Log.Fatalf("Redis ping failed: %v", err)
	}
	utils.Log.Infof("Redis connected: %s", redisURL)
	return rdb
}

func setupFiberApp() *fiber.App {
	app := fiber.New(config.FiberConfig())

	// Middleware setup
	app.Use("/api/auth", middleware.LimiterConfig())
	app.Use(middleware.LoggerConfig())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.RecoverConfig())

	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	return db
}

func setupRoutes(app *fiber.App, db *gorm.DB, rdb *redis.Client) {

	// route.Routes(app, db)
	// app.Use(utils.NotFoundHandler)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "api",
			"version": os.Getenv("VERSION"),
		})
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error", "db": "unavailable", "redis": "unknown",
			})
		}
		ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cancel()

		dbOK := sqlDB.PingContext(ctx) == nil
		redisOK := rdb.Ping(ctx).Err() == nil

		status := fiber.StatusOK
		statusText := "ok"
		if !dbOK || !redisOK {
			status = fiber.StatusServiceUnavailable
			statusText = "degraded"
		}
		body := fiber.Map{
			"status": statusText,
			"db":     map[bool]string{true: "up", false: "down"}[dbOK],
			"redis":  map[bool]string{true: "up", false: "down"}[redisOK],
		}
		return c.Status(status).JSON(body)
	})

	route.Routes(app, db)
	app.Use(utils.NotFoundHandler)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	if err := app.Listen(address); err != nil {
		errs <- fmt.Errorf("error starting server: %w", err)
	}
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("Error closing database connection: %v", err)
	} else {
		utils.Log.Info("Database connection closed successfully")
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		utils.Log.Fatalf("Server error: %v", err)
	case <-quit:
		utils.Log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			utils.Log.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		utils.Log.Info("Server exiting due to context cancellation")
	}

	utils.Log.Info("Server exited")
}
