package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

var ErrNoDefaults = errors.New("no default config")

var (
	ErrMergingError = errors.New("merging with defaults error")
)

// Read Reads yaml files from configuration directory with sub folders
// as application stage and merges config files in one configuration per stage.
func Read(ctx context.Context, stage Stageable, opts ...Option) ([]byte, error) {
	mc, err := newMerger(append(opts, withStageName(stage))...)
	if err != nil {
		return nil, err
	}

	mc.logger.Debug(ctx, "Current stage", mc.stage.Name().String())
	mc.logger.Debug(ctx, "Config path", mc.cfgPath)

	if err = mc.getFileList(); err != nil {
		mc.logger.Error(ctx, "get file list error", err)

		return nil, err
	}

	// check defaults config existence. Fall down if not
	if !mc.defaultConfigExists() {
		mc.logger.Error(ctx, "defaults config is not found in file list! Fall down.", mc.fileList)

		return nil, ErrNoDefaults
	}

	mc.logger.Debug(ctx, "Existing config list", mc.fileList)

	fileListResult := make(map[StageName][]string)

	for stageName, files := range mc.fileList {
		for _, file := range files {
			var configBytes []byte

			configBytes, err = os.ReadFile(file)
			if err != nil {
				mc.logger.Error(ctx, "Config read fail! Fall down.", stageName, file)

				return nil, fmt.Errorf("config file `%s` read fail: %w", file, err)
			}

			var configFromFile map[StageName]map[string]any

			mc.logger.Debug(ctx, "file content", file, string(configBytes))

			if err = yaml.Unmarshal(configBytes, &configFromFile); err != nil {
				mc.logger.Error(ctx, "config read fail! Fall down.", stageName, file)

				return nil, fmt.Errorf("config file `%s` unmarshal fail: %w", file, err)
			}

			if _, ok := configFromFile[stageName]; !ok {
				mc.logger.Warn(ctx, "File excluded from current stage (it is not for this stage)!", file, stageName)

				continue
			}

			cc, ok := mc.getConfigForStage(stageName)
			if !ok {
				mc.setConfigForStage(stageName, configFromFile[stageName])
				cc, _ = mc.getConfigForStage(stageName)
			}

			if err = mergo.Merge(
				&cc,
				configFromFile[stageName],
			); err != nil {
				mc.logger.Error(ctx, "config merge fail! Fall down.", stageName, file)

				return nil, fmt.Errorf("merging configs[%s] with configFromFile[%s] config fail: %w", stageName, stageName, err)
			}

			mc.setConfigForStage(stageName, cc)

			fileListResult[stageName] = append(fileListResult[stageName], file)
		}
	}

	mc.logger.Debug(ctx, "Parsed config list", fileListResult)

	return mc.getResultConfigForStage(ctx)
}

func (mc *merger) getResultConfigForStage(ctx context.Context) ([]byte, error) {
	resultConfig, ok := mc.getConfigForStage(StageNameDefaults)
	if !ok {
		return nil, ErrNoDefaults
	}

	if cfg, ok := mc.getConfigForStage(mc.stage.Name()); ok {
		if err := mergo.Merge(&resultConfig, cfg, mergo.WithOverride); err != nil {
			mc.logger.Error(ctx, "merging with defaults error: %s", err)

			return nil, fmt.Errorf("%w: %w", ErrMergingError, err)
		}

		mc.logger.Debug(ctx, "Stage config is loaded and merged with `defaults`", mc.stage.Name().String())
	}

	mc.logger.Debug(ctx, "final config", resultConfig)

	return yaml.Marshal(resultConfig)
}

func (mc *merger) getConfigForStage(stage StageName) (map[string]any, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if len(mc.configs[stage]) == 0 {
		return nil, false
	}

	cfg := make(map[string]any, len(mc.configs[stage]))
	for key, val := range mc.configs[stage] {
		cfg[key] = val
	}

	return cfg, true

}

func (mc *merger) getFileList() error {
	var stageDir StageName

	return filepath.Walk(mc.cfgPath, func(path string, f os.FileInfo, err error) error {
		if mc.cfgPath == path {
			return nil
		}

		if f.IsDir() {
			if stageDir.String() == "" || f.Name() == StageNameDefaults.String() || mc.stage.Name().String() == f.Name() {
				stageDir = NewStageNameUnsafe(f.Name())

				return nil
			}

			return filepath.SkipDir
		}

		if fileIsYaml(f.Name()) && (stageDir.isDefault() || mc.isSpecifiedStage(stageDir)) {
			mc.addFileToStage(stageDir, f.Name())
		}

		return nil
	})
}

func fileIsYaml(name string) bool {
	return filepath.Ext(name) == ".yaml" || filepath.Ext(name) == ".yml"
}

func (mc *merger) isSpecifiedStage(stage StageName) bool {
	return mc.stage.Name() == stage
}

func (mc *merger) addFileToStage(stage StageName, file string) {
	fileList, ok := mc.fileList[stage]
	if !ok || len(fileList) == 0 {
		fileList = make([]string, 0, 1)
	}

	file = mc.cfgPath + "/" + stage.String() + "/" + file

	fileList = append(fileList, file)

	mc.fileList[stage] = fileList
}

func (mc *merger) defaultConfigExists() bool {
	defaultConfig, ok := mc.fileList[StageNameDefaults]
	if !ok || len(defaultConfig) == 0 {
		return false
	}

	return true
}

func (mc *merger) setConfigForStage(stageName StageName, cfg map[string]any) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.configs[stageName] = cfg
}
