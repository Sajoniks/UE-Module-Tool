package ue

import "path/filepath"

type ProjectModuleDescriptor struct {
	Name                   string       `json:"Name"`
	LoadingPhase           LoadingPhase `json:"LoadingPhase"`
	Type                   ModuleType   `json:"Type"`
	AdditionalDependencies []string     `json:"AdditionalDependencies,omitempty"`
}

type PluginDescriptor struct {
	Name                     string   `json:"Name"`
	Enabled                  bool     `json:"Enabled"`
	SupportedTargetPlatforms []string `json:"SupportedTargetPlatforms,omitempty"`
}

type ProjectFileDescriptor struct {
	ProjectPath     string `json:"-"`
	ProjectFileName string `json:"-"`
	ProjectName     string `json:"-"`
	IsPlugin        bool   `json:"-"`

	FileVersion       int                        `json:"FileVersion"`
	EngineAssociation string                     `json:"EngineAssociation"`
	Category          string                     `json:"Category,omitempty"`
	Description       string                     `json:"Description,omitempty"`
	Modules           []*ProjectModuleDescriptor `json:"Modules,omitempty"`
	Plugins           []*PluginDescriptor        `json:"Plugins,omitempty"`
	TargetPlatforms   []string                   `json:"TargetPlatforms,omitempty"`
}

func (p *ProjectFileDescriptor) Path() string {
	return filepath.Join(p.ProjectPath, p.ProjectFileName)
}

func (p *ProjectFileDescriptor) Sources() string {
	return filepath.Join(p.ProjectPath, "Source")
}

func (p *ProjectFileDescriptor) ModuleSources(mdl string) string {
	return filepath.Join(p.ProjectPath, "Source", mdl)
}

func (p *ProjectFileDescriptor) ModulePublic(mdl string) string {
	return filepath.Join(p.ProjectPath, "Source", mdl, "Public")
}

func (p *ProjectFileDescriptor) ModulePrivate(mdl string) string {
	return filepath.Join(p.ProjectPath, "Source", mdl, "Private")
}

func (p *ProjectFileDescriptor) Touch() {
	// format for ue4.27
	p.FileVersion = 3
	p.EngineAssociation = "4.27"
}
