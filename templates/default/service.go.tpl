package domains

type AppService struct{}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) GetApp() string {
	return "Hello World"
}
