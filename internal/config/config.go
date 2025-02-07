package config

import "github.com/MetsysEht/setuProject/pkg/gormDatabase"

// App contains application-specific config values

type Config struct {
	App                App
	DB                 gormDatabase.Config
	SetuGatewayService SetuGatewayService
}

type App struct {
	Env             string
	ServiceName     string
	Hostname        string
	Port            string
	ShutdownTimeout int
	ShutdownDelay   int
	GitCommitHash   string
	Key             string
	Interfaces      struct {
		Service NetworkInterfaces
	}
}

type NetworkInterfaces struct {
	GrpcServerAddress     string
	HttpServerAddress     string
	InternalServerAddress string
}
