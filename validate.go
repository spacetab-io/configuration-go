package config

type Config interface {
	ValidateAll() []error
}

type Configurator interface {
	Validate() []error
}
