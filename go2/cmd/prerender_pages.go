package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	pathpkg "path"
	"strings"
	"text/template"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zaydek/retro/errs"
)

// PrerenderedPage describes a response for a page prerendered by Node.
type PrerenderedPage struct {
	FSPath string `json:"fs_path"`
	Path   string `json:"path"`
	Head   string `json:"head"`
	Page   string `json:"page"`
}

func prerenderPages(app *RetroApp) error {
	bstr, err := ioutil.ReadFile(pathpkg.Join(app.Configuration.AssetDirectory, "index.html"))
	if err != nil {
		return errs.ReadFile(pathpkg.Join(app.Configuration.AssetDirectory, "index.html"), err)
	}

	text := string(bstr)
	if !strings.Contains(text, "{{ .Head }}") {
		return errors.New("No such template tag '{{ .Head }}'. " +
			"This is the entry point for the '<Head>' component in your page components. " +
			"Add '{{ .Head }}' to '<head>'.")
	} else if !strings.Contains(text, "{{ .Page }}") {
		return errors.New("No such template tag '{{ .Page }}'. " +
			"This is the entry point for the '<Page>' component in your page components. " +
			"Add '{{ .Page }}' to '<body>'.")
	}

	tmpl, err := template.New(pathpkg.Join(app.Configuration.AssetDirectory, "index.html")).Parse(text)
	if err != nil {
		return errs.ParseTemplate(pathpkg.Join(app.Configuration.AssetDirectory, "index.html"), err)
	}

	rawstr := `// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.

import React from "react"
import ReactDOMServer from "react-dom/server"

// Pages
` + buildRequireStmt(app.PageBasedRouter) + `

// Props
const props = require("../` + app.Configuration.CacheDirectory + `/props.js").default

async function asyncRun(requireStmtAsArray) {
	const chain = []
	for (const { fs_path, path, exports } of requireStmtAsArray) {
		const promise = new Promise(async resolve => {
			const { Head, default: Page } = exports

			// Resolve <Head {...props}>:
			let head = ""
			if (Head) {
 				head = ReactDOMServer.renderToStaticMarkup(<Head {...props[path]} />)
			}
			head = head.replace(/></g, ">\n\t\t<")
			head = head.replace(/\/>/g, " />")

			// Resolve <Page {...props}>:
			let page = '<div id="root"></div>'
			if (Page) {
				page = ReactDOMServer.renderToString(<div id="root"><Page {...props[path]} /></div>)
			}
			page += '\n\t\t<script src="/app.js"></script>'

			resolve({ fs_path, path, head, page })
		})
		chain.push(promise)
	}
	const resolvedAsArr = await Promise.all(chain)
	console.log(JSON.stringify(resolvedAsArr, null, 2))
}

asyncRun(` + buildRequireStmtAsArray(app.PageBasedRouter) + `)
`

	if err := ioutil.WriteFile(pathpkg.Join(app.Configuration.CacheDirectory, "pages.esbuild.js"), []byte(rawstr), 0644); err != nil {
		return errs.WriteFile(pathpkg.Join(app.Configuration.CacheDirectory, "pages.esbuild.js"), err)
	}

	results := api.Build(api.BuildOptions{
		Bundle: true,
		Define: map[string]string{
			"__DEV__":              fmt.Sprintf("%t", app.Configuration.Env == "development"),
			"process.env.NODE_ENV": fmt.Sprintf("%q", app.Configuration.Env),
		},
		EntryPoints: []string{pathpkg.Join(app.Configuration.CacheDirectory, "pages.esbuild.js")},
		Loader:      map[string]api.Loader{".js": api.LoaderJSX},
	})
	if len(results.Errors) > 0 {
		bstr, err := json.MarshalIndent(results.Errors, "", "\t")
		if err != nil {
			return errs.Unexpected(err)
		}
		return errors.New(string(bstr))
	}

	stdoutBuf, err := pipeToNode(results.OutputFiles[0].Contents)
	if err != nil {
		return err
	}

	var pages []PrerenderedPage
	if err := json.Unmarshal(stdoutBuf.Bytes(), &pages); err != nil {
		return errs.Unexpected(err)
	}

	// TODO: Change to sync.WaitGroup or errgroup?
	for _, each := range pages {
		var path string
		path = each.FSPath[len(app.Configuration.PagesDirectory):]  // pages/page.js -> page.js
		path = path[:len(path)-len(pathpkg.Ext(path))] + ".html"    // page.js -> page.html
		path = pathpkg.Join(app.Configuration.BuildDirectory, path) // page.html -> build/page.html
		if dir := pathpkg.Dir(path); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return errs.MkdirAll(dir, err)
			}
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, each); err != nil {
			return errs.ExecuteTemplate(path, err)
		}
		if err := ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
			return errs.WriteFile(path, err)
		}
	}
	return nil
}
