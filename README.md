Golang Microservice configuration module
----------------------------------------

[![codecov](https://codecov.io/gh/spacetab-io/configuration-go/graph/badge.svg)](https://codecov.io/gh/spacetab-io/configuration-go)

Configuration module for microservices written on Go.
Preserves [corporate standards for services configuration](https://confluence.teamc.io/pages/viewpage.action?pageId=4227704).

## Installation

Import in your configuration file

```go
package main

import (
	config "github.com/spacetab-io/configuration-go"
)

```

## Usage

Some agreements:

1. Configuration must be declared as struct and reveals yaml structure
2. Default config folder: `./configuration`. If you need to override, pass your path in `ReadConfig` function
3. Stage is passed as `config.Stageable` implementation. In example below stageEnv is used to pass stage through env variable `STAGE`.

Code example:

```go
package main

import (
	"context"
	"fmt"
	"log"

	config "github.com/spacetab-io/configuration-go"
	"github.com/spacetab-io/configuration-go/stage"
	"gopkg.in/yaml.v3"
)

// ConfigStruct is your app config structure. This must be related to yaml config file structure. 
// Everything that is not in this struct will be passed and will not be processed.
// Keep in mind that inheritance must be implemented with `struct{}`
type ConfigStruct struct {
	Log struct {
		Level  string `yaml:"level"`
		Debug  bool   `yaml:"debug"`
		Format string `yaml:"yaml"`
	} `yaml:"log"`
	Database struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		User           string `yaml:"user"`
		Pass           string `yaml:"password"`
		Name           string `yaml:"name"`
		SslMode        string `yaml:"sslMode"`
		Logs           bool   `yaml:"logs"`
		MigrateOnStart bool   `yaml:"migrateOnStart"`
		MigrationPath  string `yaml:"migrationsPath"`
	} `yaml:"db"`
	Http struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"ws"`
}

func main() {
	// config.Read receives stage as config.Stageable implementation.
	// You can use envStage to pass stage name via ENV param STAGE.
	// In NewEnvStage you can pass fallback value if STAGE param is empty.
	envStage := config.NewEnvStage("development")
	// Reading ALL config files in defaults configuration folder and recursively merge them with STAGE configs
	configBytes, err := config.Read(context.TODO(), envStage, config.WithConfigPath("./configuration"))
	if err != nil {
		log.Fatalf("config reading error: %+v", err)
	}

	var cfg ConfigStruct
	// unmarshal config into Config structure 
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		log.Fatalf("config unmarshal error: %+v", err)
	}

	fmt.Printf("config: %+v", cfg)
} 
```

## License

The MIT License

Copyright © 2024 SpaceTab.io, Inc. https://spacetab.io

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "
Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the
following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.