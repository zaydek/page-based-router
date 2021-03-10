import * as child_process from "child_process"
import * as fs from "fs"
import * as path from "path"

import { copyAll, readdirAll } from "../fsAll"

test("copyAll", async () => {
	await fs.promises.mkdir(path.join(__dirname, "foo"), { recursive: true })
	await fs.promises.mkdir(path.join(__dirname, "foo", "bar"), { recursive: true })
	await fs.promises.mkdir(path.join(__dirname, "foo", "bar", "baz"), { recursive: true })

	await fs.promises.writeFile(path.join(__dirname, "foo/a"), "")
	await fs.promises.writeFile(path.join(__dirname, "foo/bar/b"), "")
	await fs.promises.writeFile(path.join(__dirname, "foo/bar/baz/c"), "")
	await fs.promises.writeFile(path.join(__dirname, "foo/bar/baz/exclude"), "")

	await copyAll(path.join(__dirname, "foo"), path.join(__dirname, "bar"), [path.join(__dirname, "foo/bar/baz/exclude")])

	let fooSrcs = await readdirAll(path.join(__dirname, "foo"))
	fooSrcs = fooSrcs.map(src => path.relative(__dirname, src))

	let barSrcs = await readdirAll(path.join(__dirname, "bar"))
	barSrcs = barSrcs.map(src => path.relative(__dirname, src))

	// prettier-ignore
	expect(fooSrcs).toEqual([
		"foo/a",
		"foo/bar",
		"foo/bar/b",
		"foo/bar/baz",
		"foo/bar/baz/c",
		"foo/bar/baz/exclude",
	])

	// prettier-ignore
	expect(barSrcs).toEqual([
		"bar/a",
		"bar/bar",
		"bar/bar/b",
		"bar/bar/baz",
		"bar/bar/baz/c",
	])

	// // NOTE: fs.promises.unlink throws EPERM error.
	// await fs.promises.unlink(path.join(__dirname, "foo"))
	child_process.execSync(`rm -rf ${path.join(__dirname, "foo")}`)
	child_process.execSync(`rm -rf ${path.join(__dirname, "bar")}`)
})

test("readdirAll", async () => {
	await fs.promises.mkdir(path.join(__dirname, "foo"), { recursive: true })
	await fs.promises.mkdir(path.join(__dirname, "foo", "bar"), { recursive: true })
	await fs.promises.mkdir(path.join(__dirname, "foo", "bar", "baz"), { recursive: true })

	await fs.promises.writeFile(path.join(__dirname, "foo/a"), "")
	await fs.promises.writeFile(path.join(__dirname, "foo/bar/b"), "")
	await fs.promises.writeFile(path.join(__dirname, "foo/bar/baz/c"), "")
	await fs.promises.writeFile(path.join(__dirname, "foo/bar/baz/exclude"), "")

	let srcs = await readdirAll(path.join(__dirname, "foo"), [path.join(__dirname, "foo/bar/baz/exclude")])
	srcs = srcs.map(src => path.relative(__dirname, src))

	// prettier-ignore
	expect(srcs).toEqual([
		"foo/a",
		"foo/bar",
		"foo/bar/b",
		"foo/bar/baz",
		"foo/bar/baz/c",
	])

	// // NOTE: fs.promises.unlink throws EPERM error.
	// await fs.promises.unlink(path.join(__dirname, "foo"))
	child_process.execSync(`rm -rf ${path.join(__dirname, "foo")}`)
})
