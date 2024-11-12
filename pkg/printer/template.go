package printer

import (
	"fmt"
	"strings"
	"text/template"
)

var tplFuncs = template.FuncMap{
	"split": func(args ...any) []string {
		sep := args[0].(string)
		str := args[1].(string)
		ls := strings.Split(str, sep)
		cp := make([]string, 0, len(ls))
		for i, _ := range ls {
			ls[i] = strings.TrimSpace(ls[i])
			if len(ls[i]) > 0 {
				cp = append(cp, ls[i])
			}
		}
		return cp
	},
}

var globalTpl *template.Template

type templateString func() (string, string)

func fromString(name, str string) templateString {
	return func() (string, string) {
		return name, str
	}
}

func loadTemplates(fs ...templateString) (*template.Template, error) {
	var err error
	var tpl *template.Template
	for _, f := range fs {
		name, str := f()
		if tpl == nil {
			tpl, err = template.New(name).Funcs(tplFuncs).Parse(str)
		} else {
			tpl, err = tpl.New(name).Parse(str)
		}

		if err != nil {
			return nil, err
		}
	}
	return tpl, nil
}

func init() {
	var err error
	globalTpl, err = loadTemplates(
		fromString("copyright", templateCopyright),
		fromString("header", templateHeader),
		fromString("build_file", templateBuildCs),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load templates: %v", err))
	}
}

func moduleTemplate() *template.Template {
	return globalTpl.Lookup("header")
}
func moduleBuildFile() *template.Template {
	return globalTpl.Lookup("build_file")
}
