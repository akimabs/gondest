package main

import (
	"context"
	"{{ .AppName }}/domains"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.uber.org/fx"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	fxApp := fx.New(
		fx.Provide(func() *fiber.App {
			return app
		}),
		
		// Provide the domain module here
		domains.Module,

		fx.Invoke(func(appController *domains.AppController) {
			// Register routes here
			appController.RegisterRoutes(app)

			if err := app.Listen(":3000"); err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}),
	)

	if err := fxApp.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start FX app: %v", err)
	}
}
