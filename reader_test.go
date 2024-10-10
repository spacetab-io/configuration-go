package config_test

import (
	"context"
	"os"
	"testing"

	config "github.com/spacetab-io/configuration-go"
	"github.com/spacetab-io/configuration-go/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadConfigs(t *testing.T) {
	t.Parallel()

	t.Run("Success parsing common dirs and files", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("dev")
		configBytes, err := config.Read(context.TODO(), tStage, config.WithConfigPath("./config_examples/configuration"))
		require.NoError(t, err)

		type logCfg struct {
			Level  string `yaml:"level"`
			Format string `yaml:"format"`
		}

		type cfg struct {
			Log         logCfg `yaml:"log"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
			StringValue string `yaml:"string_test"`
			Debug       bool   `yaml:"debug"`
			BoolValue   bool   `yaml:"bool_test"`
		}

		var exp cfg
		require.NoError(t, yaml.Unmarshal(configBytes, &exp))

		refConfig := cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "error", Format: "text"},
			Host:        "127.0.0.1",
			Port:        "8888",
			StringValue: "",
			BoolValue:   false,
		}

		assert.EqualValues(t, refConfig, exp)
	})

	t.Run("Success parsing merging many config files in default and one file in stage", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("local")
		configBytes, err := config.Read(context.TODO(), tStage, config.WithConfigPath("./config_examples/configuration"))
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Redis struct {
				Hostname string `yaml:"hostname"`
				Password string `yaml:"password"`
				Database int    `yaml:"database"`
				Port     int    `yaml:"port"`
			} `yaml:"redis"`
			Log struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
			StringValue string `yaml:"string_test"`
			Debug       bool   `yaml:"debug"`
			BoolValue   bool   `yaml:"bool_test"`
		}

		exp := &cfg{}
		err = yaml.Unmarshal(configBytes, &exp)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			Debug: true,
			Redis: struct {
				Hostname string `yaml:"hostname"`
				Password string `yaml:"password"`
				Database int    `yaml:"database"`
				Port     int    `yaml:"port"`
			}{
				Hostname: "127.1.1.1",
				Password: "very-very-secure-password-with-a-lot-of-words-to-make-sure-that-it-length-is-more-than-100-chars-length",
				Database: 123,
				Port:     321,
			},
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "debug", Format: "русский мат"},
			Host:        "0.0.0.0",
			Port:        "6666",
			StringValue: "not a simple string",
			BoolValue:   false,
		}

		assert.EqualValues(t, refConfig, exp)
	})

	t.Run("Success parsing common dirs and files with different stages", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("prod")
		configBytes, err := config.Read(context.TODO(), tStage, config.WithConfigPath("./config_examples/configuration"))
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

		exp := &cfg{}
		err = yaml.Unmarshal(configBytes, &exp)
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

		assert.EqualValues(t, refConfig, exp)
	})

	t.Run("Success parsing complex dirs and files", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("development")
		configBytes, err := config.Read(context.TODO(), tStage, config.WithConfigPath("./config_examples/configuration2"))
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type hbParams struct {
			AreaMapping map[string]string `yaml:"area_mapping"`
			URL         string            `yaml:"url"`
			Username    string            `yaml:"username"`
			Password    string            `yaml:"password"`
		}

		type cfg struct {
			HotelbookParams hbParams `yaml:"hotelbook_params"`
			Logging         string   `yaml:"logging"`
			DefaultList     []string `yaml:"default_list"`
			Databases       struct {
				Redis struct {
					Master struct {
						Username string `yaml:"username"`
						Password string `yaml:"password"`
					} `yaml:"master"`
				} `yaml:"redis"`
			} `yaml:"databases"`
		}

		exp := &cfg{}
		err = yaml.Unmarshal(configBytes, &exp)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			HotelbookParams: hbParams{
				AreaMapping: map[string]string{"KRK": "Krakow", "MSK": "Moscow", "CHB": "Челябинск"},
				URL:         "https://hotelbook.com/xml_endpoint",
				Username:    "TESt_USERNAME",
				Password:    "PASSWORD",
			},
			DefaultList: []string{"bar", "baz"},
			Logging:     "info",
			Databases: struct {
				Redis struct {
					Master struct {
						Username string `yaml:"username"`
						Password string `yaml:"password"`
					} `yaml:"master"`
				} `yaml:"redis"`
			}{Redis: struct {
				Master struct {
					Username string `yaml:"username"`
					Password string `yaml:"password"`
				} `yaml:"master"`
			}{Master: struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{Username: "R_USER", Password: "R_PASS"}}},
		}

		assert.EqualValues(t, refConfig, exp)
	})

	t.Run("Success parsing symlinked files and dirs", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("dev")
		configBytes, err := config.Read(context.TODO(), tStage, config.WithConfigPath("./config_examples/symnlinkedConfigs"))
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

		var exp cfg
		err = yaml.Unmarshal(configBytes, &exp)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := cfg{
			Debug: true,
			Log: struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			}{Level: "error", Format: "text"},
			Host: "127.0.0.1",
			Port: "8888",
		}

		assert.EqualValues(t, refConfig, exp)
	})

	if os.Getenv("IN_CONTAINER") == "true" {
		t.Run("Success parsing symlinked files and dirs in root", func(t *testing.T) {
			t.Parallel()
			tStage := tests.NewTestStage("dev")
			configBytes, err := config.Read(context.TODO(), tStage, config.WithConfigPath("/cfgs"))
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

			exp := &cfg{}
			err = yaml.Unmarshal(configBytes, &exp)
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

			assert.EqualValues(t, refConfig, exp)
		})
	}

	t.Run("Fail dir not found", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("dev")
		_, err := config.Read(context.TODO(), tStage, config.WithConfigPath(""))
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})

	t.Run("no defaults configs", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("dev")
		_, err := config.Read(context.TODO(), tStage, config.WithConfigPath("/config_examples/no_defaults"))
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})

	t.Run("merge errors", func(t *testing.T) {
		t.Parallel()
		tStage := tests.NewTestStage("dev")
		_, err := config.Read(context.TODO(), tStage, config.WithConfigPath("/config_examples/merge_error"))
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}
