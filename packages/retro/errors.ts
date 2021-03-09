import * as T from "./types"
import * as terminal from "../shared/terminal"
import * as utils from "./utils"

function accent(str: string): string {
	return str.replace(/('[^']+')/g, terminal.magenta("$1"))
}

function format(str: string): string {
	return accent(utils.detab(str).trim())
}

////////////////////////////////////////////////////////////////////////////////
// CLI

export function badCommand(run: string): string {
	return format(`
		Bad command '${run}'.

		Supported commands:

		dev     Start the dev server
		export  Export the production-ready build (SSG)
		serve   Serve the production-ready build
	`)
}

////////////////////////////////////////////////////////////////////////////////
// Document (index.html)

export function missingDocumentHeadTag(src: string): string {
	return format(`
		${src}: Add '%head%' somewhere to '<head>'.

		For example:

		${terminal.dim(`// ${src}`)}
		<!DOCTYPE html>
			<head lang="en">
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				${terminal.magenta("%head%")}
				${terminal.dim("...")}
			</head>
			<body>
				${terminal.dim("...")}
			</body>
		</html>
	`)
}

export function missingDocumentAppTag(src: string): string {
	return format(`
		${src}: Add '%app%' somewhere to '<body>'.

		For example:

		${terminal.dim(`// ${src}`)}
		<!DOCTYPE html>
			<head lang="en">
				${terminal.dim("...")}
			</head>
			<body>
				${terminal.magenta("%app%")}
				${terminal.dim("...")}
			</body>
		</html>
	`)
}

////////////////////////////////////////////////////////////////////////////////
// Pages

export function pagesUseNonURICharacters(pages: string[]): string {
	return format(`
		One or more pages use non URI-safe characters:

		${pages.map(page => "- " + page).join("\n")}

		URI characters described by RFC 3986:

		2.2. Unreserved Characters

			ALPHA / DIGIT / "-" / "." / "_" / "~"

		2.3. Reserved Characters

			gen-delims = ":" / "/" / "?" / "#" / "[" / "]" /
			sub-delims = "@" / "!" / "$" / "&" / "'" / "(" / ")"
			${"\x20".repeat(11)}/ "*" / "+" / "," / ";" / "="

		${terminal.underline("https://tools.ietf.org/html/rfc3986")}
	`)
}

////////////////////////////////////////////////////////////////////////////////
// Server API

export function badStaticPageExports(src: string): string {
	return format(`
		${src}: Bad static page exports.

		Static page exports should look something like this:

		${terminal.dim(`// ${src}`)}
		export function serverProps() {
			return { ${terminal.dim("...")} }
		}

		export function Head({ path, ...props }) {
			return <title>Hello, world!</title>
		}

		export default function Route({ path, ...props }) {
			return <h1>Hello, world!</h1>
		}
	`)
}

export function badDynamicPageExports(src: string): string {
	return format(`
		${src}: Bad dynamic page exports.

		Dynamic page exports should look something like this:

		${terminal.dim(`// ${src}`)}
		export function serverPaths() {
			return [
				{ path: "/foo", props: ${terminal.dim("...")} },
				{ path: "/foo/bar", props: ${terminal.dim("...")} },
				{ path: "/foo/bar/baz", props: ${terminal.dim("...")} },
			]
		}

		export function Head({ path, ...props }) {
			return <title>Hello, world!</title>
		}

		export default function Route({ path, ...props }) {
			return <h1>Hello, world!</h1>
		}
	`)
}

export function badServerPropsReturn(src: string): string {
	return format(`
		${src}.serverProps: Bad 'serverProps' return.

		A 'serverProps' resolver should look something like this:

		${terminal.dim(`// ${src}`)}
		export function serverProps() {
			return { ${terminal.dim("...")} }
		}

		${terminal.dim(`// ${src} (asynchronous)`)}
		export async function serverProps() {
			await ${terminal.dim(`...`)}
			return { ${terminal.dim("...")} }
		}
	`)
}

export function badServerPathsReturn(src: string): string {
	return format(`
		${src}.serverPaths: Bad 'serverPaths' return.

		A 'serverPaths' resolver should look something like this:

		${terminal.dim(`// ${src}`)}
		export function serverPaths() {
			return [
				{ path: "/foo", props: ${terminal.dim("...")} },
				{ path: "/foo/bar", props: ${terminal.dim("...")} },
				{ path: "/foo/bar/baz", props: ${terminal.dim("...")} },
			]
		}

		${terminal.dim(`// ${src} (asynchronous)`)}
		export async function serverPaths() {
			await ${terminal.dim(`...`)}
			return [
				{ path: "/foo", props: ${terminal.dim("...")} },
				{ path: "/foo/bar", props: ${terminal.dim("...")} },
				{ path: "/foo/bar/baz", props: ${terminal.dim("...")} },
			]
		}
	`)
}

////////////////////////////////////////////////////////////////////////////////
// Server API (etc.)

export function repeatPath(r1: T.RouteInfo, r2: T.RouteInfo): string {
	function namespace(routeInfo: T.RouteInfo): string {
		const fn = "static" ? "serverProps" : "serverPaths"
		return `${routeInfo.src}.${fn}`
	}
	let out = ""
	out += `${namespace(r1)}: Repeat path '${r1.path}'; see ${namespace(r2)}.`
	out += r1.src === r2.src ? "" : "\n"
	out += r1.src === r2.src ? "" : `\n- ${namespace(r1)} defines '${r1.path}'`
	out += r1.src === r2.src ? "" : `\n- ${namespace(r1)} defines '${r1.path}'`
	return format(out)
}

////////////////////////////////////////////////////////////////////////////////
// Serve command

// prettier-ignore
export function serveWithoutExportDirectory(): string {
	return format("App unexported; try 'retro export && retro serve'.")
}
