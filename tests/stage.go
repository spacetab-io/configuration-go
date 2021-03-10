package tests

import (
	"github.com/spacetab-io/configuration-go/stage"
)

type TestStage struct {
	name stage.Name
}

func (ts TestStage) Get() stage.Name {
	return ts.name
}

func (ts TestStage) String() string {
	return string(ts.name)
}

func NewTestStage(stageName string) stage.Interface {
	return TestStage{name: stage.Name(stageName)}
}
