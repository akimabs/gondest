package domains

import (
	"app_gondest/utils"
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
	data, err := c.service.GetApp()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.Success("App retrieved successfully", data, fiber.StatusOK))
}
