package thedeepest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	config "github.com/spacetab-io/configuration-go"
)

func TestRelativePath(t *testing.T) {
	t.Run("Success parsing relative dirs", func(t *testing.T) {
		configBytes, err := config.ReadConfigs("../../../../test/configuration")
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

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
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

		assert.EqualValues(t, refConfig, config)
	})
}
