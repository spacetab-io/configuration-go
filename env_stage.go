package config

import (
	"os"
)

type Stageable interface {
	Name() StageName
}

type EnvStage StageName

const ENVStageKey = "STAGE"

var _ Stageable = (*EnvStage)(nil)

func NewEnvStage(fallback string) Stageable {
	value, ok := os.LookupEnv(ENVStageKey)
	if !ok {
		return EnvStage(fallback)
	}

	return EnvStage(value)
}

// Name Loads stage name with fallback to 'dev'.
func (s EnvStage) Name() StageName {
	if s == "" {
		return StageNameDefaults
	}

	return StageName(s)
}
