package config

var CurrentEnvironment Environment

type Environment string

const (
	LOCAL       Environment = "local"
	TESTING     Environment = "testing"
	DEVELOPMENT Environment = "development"
	STAGING     Environment = "staging"
	ACCEPTANCE  Environment = "acceptance"
	PRODUCTION  Environment = "production"
)
