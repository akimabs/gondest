package {{ .ModuleName }}

import (
	"github.com/gofiber/fiber/v2"
)

type {{ .ControllerName }} struct {
	service *{{ .ServiceName }}
}

func New{{ .ControllerName }}(s *{{ .ServiceName }}, app *fiber.App) *{{ .ControllerName }} {
	controller := &{{ .ControllerName }}{service: s}
	controller.RegisterRoutes(app)
	return controller
}

func (c *{{ .ControllerName }}) RegisterRoutes(app *fiber.App) {
	app.Get("/{{ .ModuleName }}s", c.Get{{ .ControllerName }})
}

func (c *{{ .ControllerName }}) Get{{ .ControllerName }}(ctx *fiber.Ctx) error {
	data, err := c.service.Get{{ .ControllerName }}()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.Error(err.Error(), fiber.StatusInternalServerError))
	}
	
	return ctx.Status(fiber.StatusOK).JSON(utils.Success("App retrieved successfully", data, fiber.StatusOK))
}
