package embedded

import "text/template"

// PackageStruct describes the struct used to execute PackageTemplate.
type PackageStruct struct {
	CreateDirectory          string
	ReactVersion             string
	ReactDOMVersion          string
	RetroReactVersion        string
	RetroReactScriptsVersion string
}

//go:embed pkg.json
var pkg string

var PackageTemplate = template.Must(template.New("package.json").Parse(pkg))
