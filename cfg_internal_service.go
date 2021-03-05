package config

import (
	"time"
)

type InternalService struct {
	Enable      bool          `yaml:"enable"`
	DebugEnable bool          `yaml:"debug"`
	GzipContent bool          `yaml:"gzip_content"`
	URL         string        `yaml:"url"`
	Version     string        `yaml:"version"`
	Timeout     time.Duration `yaml:"timeout"`
}
