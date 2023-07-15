package environment

import (
	"os"
)

var env Environment
var baseDomain string

type Environment struct {
	Environment string
}

func (e Environment) IsProduction() bool {
	return e.Environment == production
}
func (e Environment) IsStaging() bool {
	return e.Environment == staging
}
func (e Environment) IsDevelopment() bool {
	return e.Environment == development
}

const (
	production  = "PRODUCTION"
	staging     = "STAGING"
	development = "DEVELOPMENT"
)

var environments = map[string]bool{
	production:  true,
	staging:     true,
	development: true,
}

func IsProduction() bool {
	if env.Environment == "" {
		getEnvironment()
	}

	return env.IsProduction()
}

func IsStaging() bool {
	if env.Environment == "" {
		getEnvironment()
	}

	return env.IsStaging()
}

func IsDevelopment() bool {
	if env.Environment == "" {
		getEnvironment()
	}

	return env.IsDevelopment()
}

func getEnvironment() {
	fenv := os.Getenv("ENV")
	if _, ok := environments[fenv]; ok {
		env = Environment{Environment: fenv}
		return
	}

	env = Environment{Environment: development}
}

func GetEnv() Environment {
	return env
}

func SetBaseDomain(domain string) {
	baseDomain = domain
}

func GetBaseDomain() string {
	return baseDomain
}
