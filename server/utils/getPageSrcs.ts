import fs from "fs"

export default function getPageSrcs() {
	// prettier-ignore
	const srcs = fs.readdirSync("pages").filter((each: string) => {
		const ok = (
			!each.startsWith("_document") &&
			each.endsWith(".tsx")
		)
		return ok
	})
	return srcs
}