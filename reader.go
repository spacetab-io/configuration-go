package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

const defaultStage = "defaults"

// ReadConfigs Reads yaml files from configuration directory with sub folders
// as application stage and merges config files in one configuration per stage
func ReadConfigs(cfgPath string) ([]byte, error) {
	if cfgPath == "" {
		cfgPath = "./configuration"
	}

	cfgPath = strings.TrimRight(cfgPath, "/")
	iSay("Config path: `%v`", cfgPath)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config path stat error: %v", err)
	}

	stage := getStage()

	var (
		fileList = map[string][]string{}
		stageDir string
	)

	_ = filepath.Walk(cfgPath, func(path string, f os.FileInfo, err error) error {
		if cfgPath == path {
			return nil
		}

		if stageDir == "" || f.Name() == defaultStage || f.Name() == stage {
			stageDir = f.Name()
		}

		if f.IsDir() {
			return nil
		}

		if filepath.Ext(f.Name()) == ".yaml" && (stageDir == defaultStage || stageDir == stage) {
			fileList[stageDir] = append(fileList[stageDir], f.Name())
		}

		return nil
	})

	// check defaults config existence. Fall down if not
	if _, ok := fileList[defaultStage]; !ok || len(fileList[defaultStage]) == 0 {
		iSay("defaults config is not found in file list `%+v`! Fall down.", fileList)
		return nil, fmt.Errorf("no default config")
	}

	iSay("Existing config list: %+v", fileList)

	fileListResult := make(map[string][]string)
	configs := make(map[string]map[string]interface{})

	for folder, files := range fileList {
		for _, file := range files {
			fullFilePath := cfgPath + "/" + folder + "/" + file
			configBytes, _ := ioutil.ReadFile(fullFilePath)

			var configFromFile map[string]map[string]interface{}

			if err := yaml.Unmarshal(configBytes, &configFromFile); err != nil {
				iSay("[config] %s %s config read fail! Fall down.", folder, file)
				return nil, fmt.Errorf("config file `%s` read fail", fullFilePath)
			}

			if _, ok := configFromFile[folder]; !ok {
				iSay("File %s excluded from %s (it is not for this stage)!", file, folder)
				continue
			}

			if _, ok := configs[folder]; !ok {
				configs[folder] = configFromFile[folder]
			}

			cc := configs[folder]
			_ = mergo.Merge(&cc, configFromFile[folder], mergo.WithOverwriteWithEmptyValue)

			configs[folder] = cc

			fileListResult[folder] = append(fileListResult[folder], file)
		}
	}

	iSay("Parsed config list: `%+v`", fileListResult)

	config := configs[defaultStage]

	if c, ok := configs[stage]; ok {
		if err := mergo.Merge(&config, c, mergo.WithOverwriteWithEmptyValue); err == nil {
			iSay("Stage `%s` config is loaded and merged with `defaults`", stage)
		}
	}

	return yaml.Marshal(config)
}

// iSay Logs in stdout when quiet mode is off
func iSay(pattern string, args ...interface{}) {
	// if quietMode == false {
	log.Printf("[config] "+pattern, args...)
	// }
}

// getStage Load configuration for stage with fallback to 'dev'
func getStage() (stage string) {
	stage = GetEnv("STAGE", "development")
	iSay("Current stage: `%s`", stage)

	return
}

// GetEnv Getting var from ENV with fallback param on empty
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
