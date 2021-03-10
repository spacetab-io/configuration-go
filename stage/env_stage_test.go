package stage_test

import (
	"testing"

	"github.com/spacetab-io/configuration-go/stage"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvStage(t *testing.T) {
	t.Parallel()

	envStage := stage.NewEnvStage("dev")

	assert.Equal(t, stage.EnvStage{Name: "dev"}, envStage)
}

func TestEnvStage_Get(t *testing.T) {
	t.Parallel()

	envStage := stage.NewEnvStage("dev")

	assert.Equal(t, stage.Name("dev"), envStage.Get())
}

func TestEnvStage_String(t *testing.T) {
	t.Parallel()

	envStage := stage.NewEnvStage("dev")

	assert.Equal(t, "dev", envStage.String())
}
