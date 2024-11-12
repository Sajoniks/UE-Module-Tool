package ue

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type LoadingPhase int

const (
	LoadingPhasePreDefault LoadingPhase = iota
	LoadingPhaseDefault
	LoadingPhasePostEngineInit
)

const (
	lpPreDefault     = "PreDefault"
	lpDefault        = "Default"
	lpPostEngineInit = "PostEngineInit"
)

func (lp LoadingPhase) String() string {
	switch lp {
	case LoadingPhaseDefault:
		return lpDefault
	case LoadingPhasePreDefault:
		return lpPreDefault
	case LoadingPhasePostEngineInit:
		return lpPostEngineInit
	}
	return ""
}

func stringToLp(str string) LoadingPhase {
	switch str {
	case lpPreDefault:
		return LoadingPhasePreDefault
	case lpDefault:
		return LoadingPhaseDefault
	case lpPostEngineInit:
		return LoadingPhasePostEngineInit
	}
	return -1
}

func (lp *LoadingPhase) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	*lp = stringToLp(str)
	return nil
}

func (lp *LoadingPhase) MarshalJSON() ([]byte, error) {
	return json.Marshal(lp.String())
}

func (lp *LoadingPhase) UnmarshalYAML(value *yaml.Node) error {
	var str string
	err := value.Decode(&str)
	if err != nil {
		return err
	}
	*lp = stringToLp(str)
	return nil
}

func (lp *LoadingPhase) MarshalYAML() (interface{}, error) {
	return lp.String(), nil
}
