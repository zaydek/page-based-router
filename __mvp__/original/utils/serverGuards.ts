import conf from "../conf"
import fs from "fs"
import { detab } from "../../utils"

// Server guards; these must run before server commands are run.
export default function serverGuards() {
	if (!fs.existsSync(conf.CACHE_DIR)) {
		fs.mkdirSync(conf.CACHE_DIR)
	}
	// Guarantee `prerender-html` can run before **or** after `prerender-props`:
	if (!fs.existsSync(conf.CACHE_DIR + "/pageProps.js")) {
		const out = detab(`
			// THIS FILE IS AUTOGENERATED.
			// THESE AREN’T THE FILES YOU’RE LOOKING FOR. MOVE ALONG.

			module.exports = {}`)
		fs.writeFileSync(conf.CACHE_DIR + "/pageProps.js", out + "\n")
	}
	if (!fs.existsSync(conf.BUILD_DIR)) {
		fs.mkdirSync(conf.BUILD_DIR)
	}
}
