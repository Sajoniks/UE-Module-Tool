package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sajoniks/ue-tools/module-tool/pkg/config"
	"github.com/sajoniks/ue-tools/module-tool/pkg/printer"
	"github.com/sajoniks/ue-tools/module-tool/pkg/ue"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func readProjectDescriptor(reader io.Reader, plugin bool) (*ue.ProjectFileDescriptor, error) {
	desc := new(ue.ProjectFileDescriptor)
	err := json.NewDecoder(reader).Decode(&desc)
	if err != nil {
		return nil, err
	}
	desc.IsPlugin = plugin
	desc.Touch()
	return desc, nil
}

func findProjectFile(dirPath string) (fs.File, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var ff os.DirEntry
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if info.IsDir() || len(info.Name()) == 0 {
			continue
		}

		ext := filepath.Ext(file.Name())
		switch ext {
		case ".uproject":
			fallthrough
		case ".uplugin":
			ff = file
		}
		if ff != nil {
			break
		}
	}

	if ff == nil {
		return nil, fmt.Errorf("project file was not found in the directory")
	}

	p := path.Join(dirPath, ff.Name())
	return os.Open(p)
}

type iOperation interface {
	Do() error
	Undo()
}

type operationStack struct {
	stack  []iOperation
	failAt int
}

func (o *operationStack) run() error {
	if o.failAt > -1 {
		return errors.New("already failed")
	}
	for ix, op := range o.stack {
		err := op.Do()
		if err != nil {
			o.failAt = ix
			return err
		}
	}
	o.failAt = -1
	return nil
}

func (o *operationStack) rollback() {
	for i := o.failAt; i >= 0; i-- {
		o.stack[i].Undo()
	}
	o.failAt = -1
}

func (o *operationStack) tryRun() error {
	err := o.run()
	if err != nil {
		o.rollback()
		return err
	}
	return nil
}

func newOperationStack(ops ...iOperation) operationStack {
	return operationStack{
		stack: ops,
	}
}

type writeProjectFileOperation struct {
	openFileHandler bool
	writtenFile     bool
	projectFile     *ue.ProjectFileDescriptor
}

func (op *writeProjectFileOperation) Do() error {
	f, err := os.Create(op.projectFile.Path())
	if err != nil {
		return err
	}
	defer f.Close()
	op.openFileHandler = true
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(&op.projectFile)
	if err != nil {
		return err
	}
	op.writtenFile = true
	return nil
}

func createProjectDirectories(projectFile *ue.ProjectFileDescriptor) error {
	_, err := os.Stat(projectFile.ProjectPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(projectFile.ProjectPath, 0666)
	}
	return err
}

func WriteProjectModule(projectFile *ue.ProjectFileDescriptor, module *ue.ProjectModuleDescriptor, cnf *config.AppConfig) error {
	var err error
	err = os.MkdirAll(projectFile.ModuleSources(module.Name), 0666)
	if err != nil {
		return err
	}
	err = os.MkdirAll(projectFile.ModulePublic(module.Name), 0666)
	if err != nil {
		return err
	}
	err = os.MkdirAll(projectFile.ModulePrivate(module.Name), 0666)
	if err != nil {
		return err
	}

	err = writeModuleBuildCs(projectFile, module.Name, cnf)
	if err != nil {
		return err
	}

	err = writeModuleCppHeader(projectFile, module.Name, cnf)
	if err != nil {
		return err
	}
	return nil
}

func writeProjectModules(projectFile *ue.ProjectFileDescriptor, cnf *config.AppConfig) error {
	for _, module := range projectFile.Modules {
		err := WriteProjectModule(projectFile, module, cnf)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeModuleBuildCs(projectFile *ue.ProjectFileDescriptor, moduleName string, cnf *config.AppConfig) error {
	var ctx printer.BuildFileCtx
	for i, _ := range cnf.Modules {
		if cnf.Modules[i].Name != moduleName {
			continue
		}
		ctx = printer.BuildFileCtx{
			Copyright:           cnf.Project.Copyright.Text,
			ModuleName:          moduleName,
			PublicDependencies:  cnf.Modules[i].Dependencies.Public,
			PrivateDependencies: cnf.Modules[i].Dependencies.Private,
		}
		break
	}
	if ctx.ModuleName != "" {
		p := filepath.Join(projectFile.ModuleSources(moduleName), moduleName+".Build.cs")
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer f.Close()
		return printer.PrintModuleBuildCs(ctx, f)
	} else {
		return errors.New("module not found")
	}
}

func writeModuleCppHeader(projectFile *ue.ProjectFileDescriptor, moduleName string, cnf *config.AppConfig) error {
	var ctx printer.ModuleCppHeaderCtx
	for i, _ := range cnf.Modules {
		if cnf.Modules[i].Name != moduleName {
			continue
		}
		ctx = printer.ModuleCppHeaderCtx{
			Copyright:  cnf.Project.Copyright.Text,
			ModuleName: moduleName,
		}
		break
	}
	if ctx.ModuleName != "" {
		p := filepath.Join(projectFile.ModuleSources(moduleName), "Public", moduleName+".h")
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer f.Close()
		return printer.PrintModuleCppHeader(ctx, f)
	} else {
		return errors.New("module not found")
	}
}

func WriteProjectFile(projectFile *ue.ProjectFileDescriptor, cnf *config.AppConfig) error {
	err := createProjectDirectories(projectFile)
	if err != nil {
		return err
	}

	err = writeProjectModules(projectFile, cnf)
	if err != nil {
		return err
	}
	return nil
}

func readFolderNames(dir string) ([]string, error) {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	folders := make([]string, 0, len(dirs))
	for _, curDir := range dirs {
		if curDir.IsDir() {
			folders = append(folders, curDir.Name())
		}
	}
	return folders, nil
}

func ReadPluginsList(projectFile *ue.ProjectFileDescriptor) ([]string, error) {
	pluginsDir := filepath.Join(projectFile.ProjectPath, "Plugins")
	return readFolderNames(pluginsDir)
}

func ReadModulesList(projectFile *ue.ProjectFileDescriptor) ([]string, error) {
	modulesDir := filepath.Join(projectFile.ProjectPath, "Source")
	return readFolderNames(modulesDir)
}

func ReadProjectFile(dirPath string) (*ue.ProjectFileDescriptor, error) {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}

	var projectFile fs.File
	if stat.IsDir() {
		projectFile, err = findProjectFile(dirPath)
	} else {
		projectFile, err = os.Open(dirPath)
	}

	if err != nil {
		return nil, err
	}
	defer projectFile.Close()

	stat, _ = projectFile.Stat()
	projPath, _ := filepath.Abs(dirPath)

	var name string
	var isPlugin bool

	ext := filepath.Ext(stat.Name())
	switch ext {
	case ".uplugin":
		isPlugin = true
	case ".uproject":
		isPlugin = false
	}
	name = stat.Name()
	name = name[:len(name)-len(ext)]

	desc, err := readProjectDescriptor(projectFile, isPlugin)
	if err != nil {
		return nil, err
	}
	desc.ProjectPath = projPath
	desc.ProjectFileName = stat.Name()
	desc.ProjectName = name

	// fix up plugins (in the list, but not in the filesystem)
	plugins, err := ReadPluginsList(desc)
	if err != nil {
		return nil, err
	}
	for _, pl := range plugins {
		for ix, ixPl := range desc.Plugins {
			if ixPl.Name != pl {
				continue
			}

			desc.Plugins = append(desc.Plugins[:ix], desc.Plugins[:ix+1]...)
			break
		}
	}

	return desc, nil
}
