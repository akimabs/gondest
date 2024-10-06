package {{ .ModuleName }}

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(New{{ .ServiceName }}),
	fx.Provide(New{{ .ControllerName }}),
)