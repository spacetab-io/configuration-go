package config_test

import (
	"testing"

	"github.com/spacetab-io/configuration-go"
	"github.com/stretchr/testify/assert"
)

func TestName_String(t *testing.T) {
	t.Parallel()

	envStage := config.NewEnvStage("dev")

	assert.Equal(t, "dev", envStage.Name().String())
}
