package main

import (
	"context"
	"log"
	"os"

	"{{ .AppName }}/domains"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/fx"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the port from the .env file, fallback to ":3000" if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]/${time} ${port} ${status} - ${method} ${path}\n",
	}))

	fxApp := fx.New(
		domains.Module,
		fx.Provide(func() *fiber.App {
			return app
		}),
		fx.Invoke(func(appController *domains.AppController) {
			// Register routes here
			appController.RegisterRoutes(app)

			if err := app.Listen(":" + port); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}),
	)

	if err := fxApp.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start FX app: %v", err)
	}
}
