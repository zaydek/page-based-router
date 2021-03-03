import * as errors from "./errors"
import * as log from "../lib/log"
import * as term from "../lib/term"
import * as types from "./types"
import * as utils from "./utils"

import newCLI from "./cli"
import newRuntimeFromCommand from "./runtime"
import runDev from "./run-dev"
import runExport from "./run-export"
import runServe from "./run-serve"

// space converts tabs to one space; "\x20".
function space(str: string): string {
	return str
		.split("\n")
		.map(each => {
			if (each.length === 0) return
			return each.replace("\t", "\x20") // "\t" -> "\x20"
		})
		.join("\n")
}

const usage = space(`
	${term.bold("Usage:")}

		retro dev          Start the dev server
		retro export       Export the production-ready build (SSG)
		retro serve        Serve the production-ready build

	${term.bold("retro dev")}

		Start the dev server

			--cached=...     Use cached resources (default false)
			--sourcemap=...  Add source maps (default true)
			--mode=...       Serve mode 'spa' or 'ssg' (default 'ssg') (experimental)
			--port=...       Port number (default 8000)

	${term.bold("retro export")}

		Export the production-ready build (SSG)

			--cached=...     Use cached resources (default false)
			--sourcemap=...  Add source maps (default true)

	${term.bold("retro serve")}

		Serve the production-ready build

			--mode=...       Serve mode 'spa' or 'ssg' (default 'ssg') (experimental)
			--port=...       Port number (default 8000)

	${term.bold("Repository")}

		${term.bold.underline.cyan("https://github.com/zaydek/retro")}
`)

async function main(): Promise<void> {
	const argv = process.argv
	if (process.argv0 === "node") {
		argv.shift()
	}

	let runCommand = "usage"
	if (argv.length >= 2) {
		runCommand = argv[1]!
	}

	let command: types.Command

	const cli = newCLI(...argv.slice(2))
	switch (runCommand) {
		case "version":
		case "--version":
		case "--v":
			console.log(process.env["RETRO_VERSION"] ?? "TODO")
			process.exit(0)
		case "usage":
		case "--usage":
			console.log(usage)
			process.exit(0)
		case "help":
		case "--help":
			console.log(usage)
			process.exit(0)
		case "dev":
			utils.setEnvDevelopment()
			command = cli.parseDevCommand()
			break
		case "export":
			utils.setEnvProduction()
			command = cli.parseExportCommand()
			break
		case "serve":
			utils.setEnvProduction()
			command = cli.parseServeCommand()
			break
		default:
			log.error(errors.badCLIRunCommand(runCommand))
			break
	}

	const runtime = await newRuntimeFromCommand(command!)
	switch (runtime.command.type) {
		case "dev":
			await runDev(runtime as types.Runtime<types.DevCommand>)
			break
		case "export":
			await runExport(runtime as types.Runtime<types.ExportCommand>)
			break
		case "serve":
			await runServe(runtime as types.Runtime<types.ServeCommand>)
			break
	}
}

process.on("uncaughtException", (err: Error): void => {
	process.env["STACK_TRACE"] = "true"
	err.message = `UncaughtException: ${err.message}`
	log.error(err)
})

main()
