package dev

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	p "path"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zaydek/retro/pkg/errs"
	"github.com/zaydek/retro/pkg/perm"
	"github.com/zaydek/retro/pkg/run"
)

func (r Runtime) RenderPageProps() error {
	src := p.Join(r.DirConfiguration.CacheDirectory, "pageProps.esbuild.js")
	dst := p.Join(r.DirConfiguration.CacheDirectory, "pageProps.js")

	// TODO: When esbuild adds support for dynamic imports, this can be changed to
	// a pure JavaScript implementation.
	text := `// THIS FILE IS AUTOGENERATED. DO NOT EDIT.

// Pages
` + strings.Join(requires(r.PageBasedRouter), "\n") + `

async function asyncRun(routes) {
	const chain = []
	for (const { path, exports, ...etc } of routes) {
		const promise = new Promise(async resolve => {
			const load = exports.load
			let props = {}
			if (load) {
				props = await load()
			}
			resolve({ ...etc, path, props })
		})
		chain.push(promise)
	}
	const resolvedAsArr = await Promise.all(chain)
	const resolvedAsMap = resolvedAsArr.reduce((acc, each) => {
		acc[each.path] = each.props
		return acc
	}, {})
	console.log(JSON.stringify(resolvedAsMap, null, 2))
}

asyncRun([
	` + strings.Join(exports(r.PageBasedRouter), ",\n\t") + `
])
`

	if err := ioutil.WriteFile(src, []byte(text), perm.File); err != nil {
		return errs.WriteFile(src, err)
	}

	results := api.Build(api.BuildOptions{
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
	})
	// TODO
	if len(results.Warnings) > 0 {
		return errors.New(formatEsbuildMessagesAsTermString(results.Warnings))
	} else if len(results.Errors) > 0 {
		return errors.New(formatEsbuildMessagesAsTermString(results.Errors))
	}

	// TODO
	stdout, err := run.Cmd(results.OutputFiles[0].Contents, "node")
	if err != nil {
		return errs.PipeEsbuildToNode(err)
	}

	contents := []byte(`// THIS FILE IS AUTOGENERATED. DO NOT EDIT.

export default ` + string(stdout))

	if err := ioutil.WriteFile(dst, contents, perm.File); err != nil {
		return errs.WriteFile(dst, err)
	}
	return nil
}
