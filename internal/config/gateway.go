package config

type Auth struct {
	Username string
	Password string
}

type SetuGatewayService struct {
	Mock                    bool
	ConnPoolConfig          ConnPoolConfig
	HystrixResiliencyConfig HystrixResiliencyConfig
	BaseUrl                 string
	Auth                    Auth
	ValidatePAN             Endpoint
	CreateRPD               Endpoint
	ClientID                string
	ClientSecret            string
	ProductInstanceID       string
}
