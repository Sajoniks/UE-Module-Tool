package factory

import (
	"fmt"
	"github.com/sajoniks/ue-tools/module-tool/pkg/ue"
	"path/filepath"
)

func CreateModule(projectFile *ue.ProjectFileDescriptor, moduleName string) (*ue.ProjectModuleDescriptor, error) {
	if moduleName == "" {
		return nil, fmt.Errorf("empty plugin name")
	}
	for _, mdl := range projectFile.Modules {
		if mdl.Name == moduleName {
			return nil, fmt.Errorf("module is already added to the project")
		}
	}

	mdl := &ue.ProjectModuleDescriptor{
		Name:         moduleName,
		LoadingPhase: ue.LoadingPhaseDefault,
		Type:         ue.ModuleRuntime,
	}
	projectFile.Modules = append(projectFile.Modules, mdl)
	return mdl, nil
}

func CreatePlugin(projectFile *ue.ProjectFileDescriptor, pluginName string, enable bool) (*ue.ProjectFileDescriptor, error) {
	if projectFile.IsPlugin {
		return nil, fmt.Errorf("can't create plugin for plugin")
	}
	if pluginName == "" {
		return nil, fmt.Errorf("empty plugin name")
	}
	for _, pl := range projectFile.Plugins {
		if pl.Name == pluginName {
			return nil, fmt.Errorf("plugin is already added to the project")
		}
	}

	// 1. create descriptor file (.uplugin)
	// 2. create default module
	// 3. modify project file descriptor (link new plugin)

	pluginDesc := ue.ProjectFileDescriptor{
		IsPlugin:          true,
		ProjectPath:       filepath.Join(projectFile.ProjectPath, "Plugins", pluginName),
		ProjectFileName:   pluginName + ".uplugin",
		ProjectName:       pluginName,
		EngineAssociation: projectFile.EngineAssociation,
		FileVersion:       projectFile.FileVersion,
	}
	_, err := CreateModule(&pluginDesc, pluginName)
	if err != nil {
		return nil, err
	}
	projectFile.Plugins = append(projectFile.Plugins,
		&ue.PluginDescriptor{
			Name:    pluginName,
			Enabled: enable,
		})
	return &pluginDesc, nil
}
