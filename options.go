package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Option func(*merger) error

var (
	ErrLoggerIsEmpty   = fmt.Errorf("logger is empty")
	ErrEmptyConfigPath = errors.New("empty config path")
	ErrConfigPath      = errors.New("config path error")
	ErrStageIsEmpty    = errors.New("stage data is empty")
)

func WithLogger(logger Logger) Option {
	return func(mc *merger) error {
		if logger == Logger(nil) {
			return ErrLoggerIsEmpty
		}

		mc.logger = logger

		return nil
	}
}

func withStageName(stage Stageable) Option {
	return func(mc *merger) error {
		if stage == nil {
			return ErrStageIsEmpty
		}

		mc.stage = stage

		return nil
	}
}

func WithConfigPath(cfgPath string) Option {
	return func(mc *merger) error {
		if cfgPath == "" {
			return ErrEmptyConfigPath
		}

		cfgPath = strings.TrimRight(cfgPath, "/")

		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrConfigPath, err)
		}

		mc.cfgPath = cfgPath

		return nil
	}
}
