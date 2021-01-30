package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	p "path"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zaydek/retro/cmd/errs"
	"github.com/zaydek/retro/perm"
)

func (r Runtime) prerenderProps() error {
	text := `// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.

// Pages
` + buildRequireStmt(r.Router) + `

async function asyncRun(requireStmtAsArray) {
	const chain = []
	for (const { path, exports } of requireStmtAsArray) {
		const promise = new Promise(async resolve => {
			const load = exports.load
			let props = {}
			if (load) {
				props = await load()
			}
			resolve({ path, props })
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

asyncRun(` + buildRequireStmtAsArray(r.Router) + `)
`

	if err := ioutil.WriteFile(p.Join(r.Config.CacheDirectory, "props.esbuild.js"), []byte(text), perm.File); err != nil {
		return errs.WriteFile(p.Join(r.Config.CacheDirectory, "props.esbuild.js"), err)
	}

	results := api.Build(api.BuildOptions{
		Bundle: true,
		Define: map[string]string{
			"__DEV__":              fmt.Sprintf("%t", os.Getenv("NODE_ENV") == "development"),
			"process.env.NODE_ENV": fmt.Sprintf("%q", os.Getenv("NODE_ENV")),
		},
		EntryPoints: []string{p.Join(r.Config.CacheDirectory, "props.esbuild.js")},
		Loader:      map[string]api.Loader{".js": api.LoaderJSX},
	})
	if len(results.Errors) > 0 {
		return errors.New(formatEsbuildMessagesAsTermString(results.Errors))
	}

	stdoutBuf, err := execNode(results.OutputFiles[0].Contents)
	if err != nil {
		return err
	}

	contents := []byte(`// THIS FILE IS AUTO-GENERATED. DO NOT EDIT.

export default ` + stdoutBuf.String())

	if err := ioutil.WriteFile(p.Join(r.Config.CacheDirectory, "props.js"), contents, perm.File); err != nil {
		return errs.WriteFile(p.Join(r.Config.CacheDirectory, "props.js"), err)
	}
	return nil
}
