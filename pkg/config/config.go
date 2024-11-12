package config

import (
	"errors"
	"fmt"
	"github.com/sajoniks/ue-tools/module-tool/pkg/ue"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type AppConfig struct {
	Project struct {
		Name      string `yaml:"name"`
		Copyright struct {
			Text      string `yaml:"text"`
			UseUnreal bool   `yaml:"use_unreal"`
		} `yaml:"copyright"`
		Category    string `yaml:"category"`
		Description string `yaml:"description"`
	} `yaml:"project"`

	Modules []struct {
		Name         string          `yaml:"name"`
		LoadingPhase ue.LoadingPhase `yaml:"loading_phase"`
		Type         ue.ModuleType   `yaml:"type"`

		Dependencies struct {
			Public  []string `yaml:"public"`
			Private []string `yaml:"private"`
		}
	} `yaml:"modules"`
}

func MustLoadProjectConfig(file string) *AppConfig {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cnf := new(AppConfig)
	err = yaml.NewDecoder(f).Decode(cnf)
	if err != nil {
		panic(err)
	}
	err = validateConfig(cnf)
	if err != nil {
		panic(err)
	}
	return cnf
}

func validateConfig(cnf *AppConfig) error {
	if strings.ContainsAny(cnf.Project.Name, "\n\t\r ") {
		return fmt.Errorf("invalid project name: %q", cnf.Project.Name)
	}
	if len(cnf.Modules) == 0 {
		return errors.New("want at least 1 module, but 0 was defined")
	}
	for _, mdl := range cnf.Modules {
		if mdl.Name == "" || strings.ContainsAny(mdl.Name, "\n\t\r ") {
			return fmt.Errorf("invalid module name: %q", mdl.Name)
		}
	}
	return nil
}
