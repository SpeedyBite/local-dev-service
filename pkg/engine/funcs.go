package engine

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// this function is copied from helm
func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	// Add some extra functionality
	extra := template.FuncMap{
		// This is a placeholder for the "include" function, which is
		// late-bound to a template. By declaring it here, we preserve the
		// integrity of the linter.
		"include":     func(string) string { return "not implemented" },
		"envVariable": func(k string, v string) string { return fmt.Sprintf("%s=%s", k, strings.Trim(v, "\"")) },
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}
