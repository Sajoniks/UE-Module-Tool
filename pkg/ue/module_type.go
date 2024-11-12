package ue

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type ModuleType int

const (
	ModuleRuntime ModuleType = iota
	ModuleEditor
	ModuleUncooked
)

const (
	mRuntime  = "Runtime"
	mEditor   = "Editor"
	mUncooked = "UncookedOnly"
)

func (m ModuleType) String() string {
	switch m {
	case ModuleRuntime:
		return mRuntime
	case ModuleEditor:
		return mEditor
	case ModuleUncooked:
		return mUncooked
	}
	return ""
}

func strToMt(str string) ModuleType {
	switch str {
	case mRuntime:
		return ModuleRuntime
	case mEditor:
		return ModuleEditor
	case mUncooked:
		return ModuleUncooked
	}
	return -1
}

func (m *ModuleType) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	*m = strToMt(str)
	return nil
}

func (m *ModuleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *ModuleType) UnmarshalYAML(value *yaml.Node) error {
	var str string
	err := value.Decode(&str)
	if err != nil {
		return err
	}
	*m = strToMt(str)
	return nil
}

func (m *ModuleType) MarshalYAML() (interface{}, error) {
	return m.String(), nil
}
