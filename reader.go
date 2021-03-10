package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	logs "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/spacetab-io/configuration-go/stage"
	"gopkg.in/yaml.v2"
)

var ErrNoDefaults = errors.New("no default config")

const defaultConfigPath = "./configuration"

// Read Reads yaml files from configuration directory with sub folders
// as application stage and merges config files in one configuration per stage.
func Read(stageI stage.Interface, cfgPath string, verbose bool) ([]byte, error) {
	cfgPath, err := checkConfigPath(cfgPath)
	if err != nil {
		return nil, err
	}

	if verbose {
		log("Current stage: `%s`", stageI.String())
		log("Config path: `%v`", cfgPath)
	}

	fileList := getFileList(stageI, cfgPath)

	// check defaults config existence. Fall down if not
	if _, ok := fileList[stage.Defaults]; !ok || len(fileList[stage.Defaults]) == 0 {
		log("defaults config is not found in file list `%+v`! Fall down.", fileList)

		return nil, ErrNoDefaults
	}

	if verbose {
		log("Existing config list: %+v", fileList)
	}

	fileListResult := make(map[stage.Name][]string)
	configs := make(map[stage.Name]map[string]interface{})

	for folder, files := range fileList {
		for _, file := range files {
			fullFilePath := cfgPath + "/" + folder.String() + "/" + file

			configBytes, err := ioutil.ReadFile(fullFilePath)
			if err != nil {
				log("%s %s config read fail! Fall down.", folder, file)

				return nil, fmt.Errorf("config file `%s` read fail: %w", fullFilePath, err)
			}

			var configFromFile map[stage.Name]map[string]interface{}

			if verbose {
				log("file `%s` content: \n%v", fullFilePath, string(configBytes))
			}

			if err := yaml.Unmarshal(configBytes, &configFromFile); err != nil {
				log("%s %s config read fail! Fall down.", folder, file)

				return nil, fmt.Errorf("config file `%s` unmarshal fail: %w", fullFilePath, err)
			}

			if _, ok := configFromFile[folder]; !ok {
				log("WARN! File `%s` excluded from `%s` (it is not for this stage)!", file, folder)

				continue
			}

			if _, ok := configs[folder]; !ok {
				configs[folder] = configFromFile[folder]
			}

			cc := configs[folder]

			err = mergo.Merge(&cc, configFromFile[folder], mergo.WithOverwriteWithEmptyValue)
			if err != nil {
				log("%s %s config merge fail! Fall down.", folder, file)

				return nil, fmt.Errorf("merging configs[%s] with configFromFile[%s] config fail: %w", folder, folder, err)
			}

			configs[folder] = cc

			fileListResult[folder] = append(fileListResult[folder], file)
		}
	}

	if verbose {
		log("Parsed config list: `%+v`", fileListResult)
	}

	config := configs[stage.Defaults]

	if c, ok := configs[stageI.Get()]; ok {
		if err := mergo.Merge(&config, c, mergo.WithOverwriteWithEmptyValue); err != nil {
			log("merging with defaults error")

			return nil, fmt.Errorf("merging with defaults error: %w", err)
		}

		log("Stage `%s` config is loaded and merged with `defaults`", stageI.String())
	}

	if verbose {
		log("final config:\n%+v", config)
	}

	return yaml.Marshal(config)
}

func getFileList(stageI stage.Interface, cfgPath string) map[stage.Name][]string {
	var (
		fileList = map[stage.Name][]string{}
		stageDir string
	)

	_ = filepath.Walk(cfgPath, func(path string, f os.FileInfo, err error) error {
		if cfgPath == path {
			return nil
		}

		if f.IsDir() {
			if stageDir == "" || f.Name() == stage.Defaults.String() || f.Name() == stageI.String() {
				stageDir = f.Name()

				return nil
			}

			return filepath.SkipDir
		}

		if filepath.Ext(f.Name()) == ".yaml" && (stageDir == stage.Defaults.String() || stageDir == stageI.String()) {
			fileList[stage.Name(stageDir)] = append(fileList[stage.Name(stageDir)], f.Name())
		}

		return nil
	})

	return fileList
}

func checkConfigPath(cfgPath string) (string, error) {
	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}

	cfgPath = strings.TrimRight(cfgPath, "/")

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return cfgPath, fmt.Errorf("config path error: %w", err)
	}

	return cfgPath, nil
}

// log Logs in stdout when quiet mode is off.
func log(pattern string, args ...interface{}) {
	logs.Printf("[config] "+pattern+"\n", args...)
}
