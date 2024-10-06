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
	data := c.service.Get{{ .ControllerName }}()
	return ctx.SendString(data)
}
