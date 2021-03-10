package stage_test

import (
	"testing"

	"github.com/spacetab-io/configuration-go/stage"
	"github.com/stretchr/testify/assert"
)

func TestName_String(t *testing.T) {
	t.Parallel()

	envStage := stage.NewEnvStage("dev")

	assert.Equal(t, "dev", envStage.String())
}
