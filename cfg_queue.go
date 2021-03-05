package config

import (
	"fmt"
)

type MessageQueue struct {
	Nsq NsqQueue `yaml:"nsq"`
}

type NsqQueue struct {
	Enable      bool   `yaml:"enable"`
	NsqdPort    int    `yaml:"nsqd_port"`
	LookupdPort int    `yaml:"lookupd_port"`
	NsqdHost    string `yaml:"nsqd_host"`
	LookupdHost string `yaml:"lookupd_host"`
	LogLevel    string `yaml:"log_level"`
}

func (s NsqQueue) GetLookupdPath() string {
	return fmt.Sprintf("%s:%v", s.LookupdHost, s.LookupdPort)
}

func (s NsqQueue) GetNsqdPath() string {
	return fmt.Sprintf("%s:%v", s.NsqdHost, s.NsqdPort)
}
