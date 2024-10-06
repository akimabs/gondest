package domains

import (
	"github.com/gofiber/fiber/v2"
)

type AppController struct {
	service *AppService
}

func NewAppController(s *AppService) *AppController {
	controller := &AppController{service: s}
	return controller
}

func (c *AppController) RegisterRoutes(app *fiber.App) {
	app.Get("/", c.GetApp)
}

func (c *AppController) GetApp(ctx *fiber.Ctx) error {
	data := c.service.GetApp()
	return ctx.SendString(data)
}
