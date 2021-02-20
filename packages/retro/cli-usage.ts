import * as term from "../lib/term"

export const cmds = `
retro dev     Starts the dev server
retro export  Exports the production-ready build (SSG)
retro serve   Serves the production-ready build
`.trim()

// prettier-ignore
export const body = `
  ${term.bold("Usage:")}

    retro dev          Starts the dev server
    retro export       Exports the production-ready build (SSG)
    retro serve        Serves the production-ready build

  ${term.boldGreen(">")} ${term.bold("retro dev")}

    Starts the dev server

      --cached         Use cached resources (default false)
      --source-map     Add source maps (default true)
      --port=<number>  Port number (default 8000)

  ${term.boldGreen(">")} ${term.bold("retro export")}

    Exports the production-ready build (SSG)

      --cached         Use cached resources (default false)
      --source-map     Add source maps (default true)

  ${term.boldGreen(">")} ${term.bold("retro serve")}

    Serves the production-ready build

      --port=<number>  Port number (default 8000)

  ${term.bold("Repository:")}

    ` + term.underline("https://github.com/zaydek/retro") + `
`