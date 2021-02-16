package retro

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	p "path"
	"regexp"
	"strings"
	"text/template"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zaydek/retro/pkg/perm"
)

// https://regex101.com/r/87JqEb/1
var appRe = regexp.MustCompile(`^app(?:\.[0-9a-f]{16})?(?:\.map)?\.js$`)

func (r Runtime) RenderApp() error {
	src := p.Join(r.DirConfiguration.CacheDirectory, "app.esbuild.js")
	dst := p.Join(r.DirConfiguration.BuildDirectory, "app.js")

	fs, err := os.ReadDir(r.DirConfiguration.BuildDirectory)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if appRe.MatchString(f.Name()) {
			path := p.Join(r.DirConfiguration.BuildDirectory, f.Name())
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}

	// TODO: When esbuild adds support for dynamic imports, this can be changed to
	// a pure JavaScript implementation.
	text := `// THIS FILE IS AUTOGENERATED. DO NOT EDIT.

import React from "react"
import ReactDOM from "react-dom"

import { BrowserRouter, Route, Router } from "@zaydek/retro-router"

// Pages
` + strings.Join(requireDefaultStmtArray(r.PageBasedRouter), "\n") + `

// Page props
const pageProps = require("./pageProps.js").default

export default function RoutedApp() {
	return (
		<BrowserRouter>
			<Router>
			{{ range $each := .PageBasedRouter }}
				<Route path="{{ $each.Path }}">
					<{{ $each.Component }} {...pageProps["{{ $each.Path }}"]} />
				</Route>
			{{ end }}
			</Router>
		</BrowserRouter>
	)
}

ReactDOM.hydrate(
	<RoutedApp />,
	document.getElementById("root"),
)
`

	var buf bytes.Buffer
	tmpl, err := template.New(src).Parse(text)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buf, r); err != nil {
		return err
	}

	if err := ioutil.WriteFile(src, buf.Bytes(), perm.File); err != nil {
		return err
	}

	result := api.Build(api.BuildOptions{
		Bundle: true,
		Define: map[string]string{
			"__DEV__":              fmt.Sprintf("%t", os.Getenv("NODE_ENV") == "development"),
			"process.env.NODE_ENV": fmt.Sprintf("%q", os.Getenv("NODE_ENV")),
		},
		EntryPoints: []string{src},
		Loader: map[string]api.Loader{
			".js": api.LoaderJSX,
			".ts": api.LoaderTSX,
		},
		MinifyIdentifiers: os.Getenv("NODE_ENV") == "production",
		MinifySyntax:      os.Getenv("NODE_ENV") == "production",
		MinifyWhitespace:  os.Getenv("NODE_ENV") == "production",
		Outfile:           dst,
		Sourcemap:         r.getSourceMap(),
		Write:             true,
	})
	// TODO
	if len(result.Errors) > 0 {
		return errors.New(formatEsbuildMessagesAsTermString(result.Errors))
	} else if len(result.Warnings) > 0 {
		return errors.New(formatEsbuildMessagesAsTermString(result.Warnings))
	}
	return nil
}