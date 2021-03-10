package stage

import (
	"os"
)

type EnvStage struct {
	Name Name
}

func NewEnvStage(fallback string) Interface {
	return EnvStage{Name: Name(getEnv("STAGE", fallback))}
}

// Get Load configuration for stage with fallback to 'dev'.
func (s EnvStage) Get() Name {
	return s.Name
}

// Get Load configuration for stage with fallback to 'dev'.
func (s EnvStage) String() string {
	return s.Name.String()
}

// getEnv Getting var from ENV with fallback param on empty.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
