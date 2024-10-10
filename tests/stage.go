package tests

import (
	"github.com/spacetab-io/configuration-go"
)

type TestStage struct {
	name config.StageName
}

func (ts TestStage) Name() config.StageName {
	return ts.name
}

func NewTestStage(stageName string) config.Stageable {
	return TestStage{name: config.StageName(stageName)}
}
