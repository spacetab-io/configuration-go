package config

type StageName string

const StageNameDefaults StageName = "defaults"

func NewStageNameUnsafe(name string) StageName {
	return StageName(name)
}

func (sn StageName) String() string {
	return string(sn)
}

func (sn StageName) isDefault() bool {
	return sn == StageNameDefaults
}
