package config_test

import (
	"os"
	"testing"

	"github.com/spacetab-io/configuration-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvStage(t *testing.T) {
	envStage := config.NewEnvStage("dev")
	assert.Equal(t, config.EnvStage("dev"), envStage)

	require.NoError(t, os.Setenv(config.ENVStageKey, "prod"))
	envStage = config.NewEnvStage("dev")
	assert.Equal(t, config.EnvStage("prod"), envStage)
	require.NoError(t, os.Unsetenv(config.ENVStageKey))
}

func TestEnvStage_Get(t *testing.T) {
	t.Parallel()

	envStage := config.NewEnvStage("dev")

	assert.Equal(t, config.StageName("dev"), envStage.Name())
}

func TestEnvStage_Name(t *testing.T) {
	t.Parallel()

	envStage := config.NewEnvStage("dev")

	assert.Equal(t, "dev", envStage.Name().String())

	envStage = config.EnvStage("")
	assert.Equal(t, config.StageNameDefaults, envStage.Name())
}
