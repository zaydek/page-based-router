import * as errors from "./errors"
import * as fsp from "fs/promises"
import * as log from "../shared/log"
import * as pages from "./pages"
import * as path from "path"
import * as router from "./router"
import * as t from "./types"
import * as utils from "./utils"

export default async function newRuntimeFromCommand(command: t.Command): Promise<t.Runtime<typeof command>> {
	const runtime: t.Runtime<typeof command> = {
		command,
		directories: {
			publicDirectory: "www",
			srcPagesDirectory: "src/pages",
			cacheDirectory: "__cache__",
			exportDirectory: "__export__",
		},
		template: "",
		pageInfos: [],
		router: {},

		// runServerGuards runs server guards.
		async runServerGuards(): Promise<void> {
			const dirs = [
				runtime.directories.publicDirectory,
				runtime.directories.srcPagesDirectory,
				runtime.directories.cacheDirectory,
				runtime.directories.exportDirectory,
			]

			for (const dir of dirs) {
				try {
					await fsp.stat(dir)
				} catch (error) {
					fsp.mkdir(dir, { recursive: true })
				}
			}

			const src = path.join(runtime.directories.publicDirectory, "index.html")

			try {
				fsp.stat(src)
			} catch (error) {
				await fsp.writeFile(
					src,
					utils.detab(`
						<!DOCTYPE html>
						<html lang="en">
							<head>
								<meta charset="utf-8" />
								<meta name="viewport" content="width=device-width, initial-scale=1" />
								%head%
							</head>
							<body>
								%page%
							</body>
						</html>
					`),
				)
			}

			const buffer = await fsp.readFile(src)
			const str = buffer.toString()

			if (!str.includes("%head")) {
				log.error(errors.missingDocumentHeadTag(src))
			} else if (!str.includes("%page")) {
				log.error(errors.missingDocumentPageTag(src))
			}
		},

		// purge purges __cache__ and __export__.
		async purge(): Promise<void> {
			const dirs = runtime.directories
			await fsp.rmdir(dirs.cacheDirectory, { recursive: true })
			await fsp.rmdir(dirs.exportDirectory, { recursive: true })

			// await this.runServerGuards()
			const excludes = [path.join(dirs.publicDirectory, "index.html")]

			// TODO: Do we need this?
			await fsp.mkdir(path.join(dirs.exportDirectory, dirs.publicDirectory), { recursive: true })
			await utils.copyAll(dirs.publicDirectory, path.join(dirs.exportDirectory, dirs.publicDirectory), excludes)
		},

		// resolveDocument resolves and or refreshes this.document.
		async resolveDocument(): Promise<void> {
			const src = path.join(this.directories.publicDirectory, "index.html")
			const buffer = await fsp.readFile(src)
			const str = buffer.toString()
			this.template = str
		},

		// resolvePages resolves and or refreshes this.pages.
		async resolvePages(): Promise<void> {
			this.pageInfos = await pages.newFromDirectories(this.directories)
		},

		// resolveRouter resolves and or refreshes this.router.
		async resolveRouter(): Promise<void> {
			this.router = await router.newFromRuntime(this)
		},
	}

	async function start(): Promise<void> {
		if (runtime.command.type === "serve") {
			// No-op
			return
		}
		await runtime.runServerGuards()
		await runtime.purge() // TODO
		await runtime.resolveDocument()
		await runtime.resolvePages()
		await runtime.resolveRouter()
	}

	await start()
	return runtime
}