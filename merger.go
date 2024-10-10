package config

import "sync"

const defaultConfigPath = "./configuration"

type merger struct {
	logger   Logger
	stage    Stageable
	cfgPath  string
	mu       sync.RWMutex
	configs  map[StageName]map[string]any
	fileList map[StageName][]string
}

func newMerger(opts ...Option) (*merger, error) {
	mc := merger{
		logger:   noOpLogger{},
		mu:       sync.RWMutex{},
		stage:    EnvStage(StageNameDefaults),
		cfgPath:  defaultConfigPath,
		fileList: make(map[StageName][]string),
		configs:  make(map[StageName]map[string]any),
	}

	for _, opt := range opts {
		if err := opt(&mc); err != nil {
			return nil, err
		}
	}

	return &mc, nil
}
