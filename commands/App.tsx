import fs from "fs"
import path from "path"
import { routeInfo } from "../Router/parts"

// Gets pages as names.
function getPageSrcs() {
	const paths = fs.readdirSync("./router")

	// prettier-ignore
	const srcs = paths.filter(each => {
		const ok = (
			fs.statSync("./router/" + each).isFile() &&
			path.parse("./router/" + each).ext === ".tsx" &&
			!each.startsWith("_") &&
			each !== "App.tsx" // TODO
		)
		return ok
	})
	const pages = srcs.map(each => path.parse(each).name)
	return pages
}

const pages = getPageSrcs()
const routeInfos = pages.map(each => routeInfo("/" + each))

// TODO: Add support for props?
function run() {
	// prettier-ignore
	fs.writeFileSync("router/App.cache.js", `
import App from "./_app"
import React from "react"
import ReactDOM from "react-dom"
import { Route, Router } from "./Router"

${routeInfos.map(each => `import ${each!.component} from ${JSON.stringify("." + each!.page)}`).join("\n")}

export default function App() {
	return (
		<Router>
			${routeInfos.map(each => `
			<Route page=${JSON.stringify(each!.page)}>
				<App>
					<${each!.component} />
				</App>
			</Route>
`).join("")}
		</Router>
	)
}

ReactDOM.render(
	<App />,
	document.getElementById("root"),
)
`.trimStart(),
	)
}

;(() => {
	run()
})()
