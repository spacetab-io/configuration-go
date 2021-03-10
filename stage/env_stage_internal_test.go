package stage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint: paralleltest // there is no parallel with env
func TestGetEnv(t *testing.T) {
	t.Run("get env key value", func(t *testing.T) {
		_ = os.Setenv("KEY", "VALUE")
		val := getEnv("KEY", "")
		if !assert.Equal(t, "VALUE", val) {
			t.FailNow()
		}
	})
	t.Run("get env key value fallback", func(t *testing.T) {
		_ = os.Setenv("KEY", "VALUE")
		val := getEnv("KEY2", "")
		if !assert.Equal(t, "", val) {
			t.FailNow()
		}
	})
}
