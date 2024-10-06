package {{ .ModuleName }}

type {{ .ServiceName }} struct{}

func New{{ .ServiceName }}() *{{ .ServiceName }} {
	return &{{ .ServiceName }}{}
}

func (s *{{ .ServiceName }}) Get{{ .ControllerName }}() (string, error) {
	return "Hi {{ .ServiceName }}", nil
}
