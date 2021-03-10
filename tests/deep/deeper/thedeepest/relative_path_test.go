package thedeepest_test

import (
	"testing"

	config "github.com/spacetab-io/configuration-go"
	"github.com/spacetab-io/configuration-go/tests"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRelativePath(t *testing.T) {
	t.Parallel()
	t.Run("Success parsing relative dirs", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("dev")
		configBytes, err := config.Read(tStage, "../../../../config_examples/configuration", false)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Log   struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		}

		conf := &cfg{}
		err = yaml.Unmarshal(configBytes, &conf)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "error", Format: "text"},
			Host: "127.0.0.1",
			Port: "8888",
		}

		assert.EqualValues(t, refConfig, conf)
	})
}
