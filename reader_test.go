package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestReadConfigs(t *testing.T) {
	t.Run("Success parsing common dirs and files", func(t *testing.T) {
		os.Setenv("STAGE", "dev")
		configBytes, err := ReadConfigs("./config_examples/configuration")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type cfg struct {
			Debug bool `yaml:"debug"`
			Log   struct {
				Level  string `yaml:"level"`
				Format string `yaml:"format"`
			} `yaml:"log"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
			StringValue string `yaml:"string_value"`
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
			Host:        "127.0.0.1",
			Port:        "8888",
			StringValue: "",
		}

		assert.EqualValues(t, refConfig, config)
	})
	t.Run("Success parsing common dirs and files with different stages", func(t *testing.T) {
		os.Setenv("STAGE", "prod")
		configBytes, err := ReadConfigs("./config_examples/configuration")
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
	t.Run("Success parsing complex dirs and files", func(t *testing.T) {
		os.Setenv("STAGE", "development")
		configBytes, err := ReadConfigs("./config_examples/configuration2")
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		type hbParams struct {
			AreaMapping map[string]string `yaml:"area_mapping"`
			Url         string            `yaml:"url"`
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

		config := &cfg{}
		err = yaml.Unmarshal(configBytes, &config)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		refConfig := &cfg{
			HotelbookParams: hbParams{
				AreaMapping: map[string]string{"KRK": "Krakow", "MSK": "Moscow", "CHB": "Челябинск"},
				Url:         "https://hotelbook.com/xml_endpoint",
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

		assert.EqualValues(t, refConfig, config)
	})

	t.Run("Success parsing symlinked files and dirs", func(t *testing.T) {
		os.Setenv("STAGE", "dev")
		configBytes, err := ReadConfigs("./config_examples/symnlinkedConfigs")
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

	if GetEnv("IN_CONTAINER", "") == "true" {
		t.Run("Success parsing symlinked files and dirs in root", func(t *testing.T) {
			os.Setenv("STAGE", "dev")
			configBytes, err := ReadConfigs("/cfgs")
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

	t.Run("Fail dir not found", func(t *testing.T) {
		_, err := ReadConfigs("")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
	t.Run("no defaults configs", func(t *testing.T) {
		_, err := ReadConfigs("./config_examples/no_defaults")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
	t.Run("merge errors", func(t *testing.T) {
		_, err := ReadConfigs("./config_examples/merge_error")
		if !assert.Error(t, err) {
			t.FailNow()
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("get env key value", func(t *testing.T) {
		os.Setenv("KEY", "VALUE")
		val := GetEnv("KEY", "")
		if !assert.Equal(t, "VALUE", val) {
			t.FailNow()
		}
	})
	t.Run("get env key value fallback", func(t *testing.T) {
		os.Setenv("KEY", "VALUE")
		val := GetEnv("KEY2", "")
		if !assert.Equal(t, "", val) {
			t.FailNow()
		}
	})
}
