package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sajoniks/ue-tools/module-tool/pkg/config"
	"github.com/sajoniks/ue-tools/module-tool/pkg/factory"
	"github.com/sajoniks/ue-tools/module-tool/pkg/parse"
	"os"
)

const App = "module-tool"

func ModuleHandler(args []string) {
	cmd, subArgs := args[0], args[1:]
	switch cmd {
	case "create":
		CreateModule(subArgs)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: - %q", cmd)
		os.Exit(-1)
	}
}

func CreateModule(args []string) {
	fs := flag.NewFlagSet("create module", flag.ExitOnError)

	var (
		cnfFilePath     = fs.String("config", "", "config file to read the plugin data from")
		projectFilePath = fs.String("project", "", "path to the .uproject or .uplugin file, or directory with this file")
	)

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	cnf := config.MustLoadProjectConfig(*cnfFilePath)
	projectFile, err := parse.ReadProjectFile(*projectFilePath)
	if err != nil {
		panic(err)
	}

	for i, _ := range cnf.Modules {
		module, createErr := factory.CreateModule(projectFile, cnf.Modules[i].Name)
		if createErr != nil {
			break
		}
		createErr = parse.WriteProjectModule(projectFile, module, cnf)
	}
	if err != nil {
		panic(err)
	}
}

func PluginHandler(args []string) {
	cmd, subArgs := args[0], args[1:]
	switch cmd {
	case "create":
		CreatePlugin(subArgs)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: - %q", cmd)
		os.Exit(-1)
	}
}

func CreatePlugin(args []string) {
	fs := flag.NewFlagSet("create plugin", flag.ExitOnError)

	var (
		cnfFilePath     = fs.String("config", "", "config file to read the plugin data from")
		projectFilePath = fs.String("project", "", "path to the .uproject or .uplugin file, or directory with this file")
	)

	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	cnf := config.MustLoadProjectConfig(*cnfFilePath)
	projectFile, err := parse.ReadProjectFile(*projectFilePath)
	if err != nil {
		panic(err)
	}

	plugin, err := factory.CreatePlugin(projectFile, cnf.Project.Name, true)
	if err != nil {
		panic(err)
	}

	err = parse.WriteProjectFile(plugin, cnf)
	if err != nil {
		panic(err)
	}
}

func main() {
	d, err := os.Getwd()
	if err != nil {
		panic(errors.New("insufficient privileges"))
	}

	fmt.Printf("Running %s in %s\n", App, d)
	flag.Parse()

	cmd, subArgs := flag.Args()[0], flag.Args()[1:]
	switch cmd {
	case "plugin":
		PluginHandler(subArgs)
	case "module":
		ModuleHandler(subArgs)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: - %q", cmd)
		os.Exit(-1)
	}
}
