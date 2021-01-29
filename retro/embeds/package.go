package embeds

import (
	_ "embed"

	"text/template"
)

type PackageDot struct {
	RepoName           string
	RetroVersion       string
	RetroRouterVersion string
	ReactVersion       string
	ReactDOMVersion    string
}

//go:embed package.json
var package_ string

var PackageTemplate = template.Must(template.New("package.json").Parse(package_))
