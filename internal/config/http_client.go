package config

import (
	"time"
)

// ConfigReader interface exposes the methods to read client config
type ConnPoolConfig struct {
	// in milliseconds
	Timeout              int
	KeepAliveTimeout     int
	MaxIdleConnections   int
	SkipCertVerification bool
}

type HystrixResiliencyConfig struct {
	MaxConcurrentRequests int
	// RequestVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	// Default is 20
	RequestVolumeThreshold int
	// CircuitBreakerSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	// Default is 5000
	CircuitBreakerSleepWindow int
	// ErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	// Default is 50
	ErrorPercentThreshold int
	// CircuitBreakerTimeout is how long to wait for command to complete, in milliseconds
	// Default is 1000
	CircuitBreakerTimeout int
}

type Endpoint struct {
	Path    string
	Method  string
	Timeout time.Duration
	Headers map[string]string `mapstructure:"headers"`
}
