package {{ .ModuleName }}

type {{ .ServiceName }} struct{}

func New{{ .ServiceName }}() *{{ .ServiceName }} {
	return &{{ .ServiceName }}{}
}

func (s *{{ .ServiceName }}) Get{{ .ControllerName }}() string {
	return "{{ .ControllerName }} Data"
}
