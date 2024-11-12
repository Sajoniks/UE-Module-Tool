package printer

import (
	"io"
)

type ModuleCppHeaderCtx struct {
	Copyright    string
	ModuleName   string
	IsGameModule bool
}

type BuildFileCtx struct {
	Copyright           string
	ModuleName          string
	PublicDependencies  []string
	PrivateDependencies []string
}

func PrintModuleCppHeader(ctx ModuleCppHeaderCtx, w io.Writer) error {
	tpl := moduleTemplate()
	if ctx.Copyright == "" {
		ctx.Copyright = defaultCopyright
	}
	err := tpl.Execute(w, &ctx)
	if err != nil {
		return err
	}
	return nil
}

func PrintModuleBuildCs(ctx BuildFileCtx, w io.Writer) error {
	tpl := moduleBuildFile()
	if ctx.Copyright == "" {
		ctx.Copyright = defaultCopyright
	}
	err := tpl.Execute(w, &ctx)
	if err != nil {
		return err
	}
	return nil
}
