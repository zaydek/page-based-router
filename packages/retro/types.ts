export interface DevCommand {
	type: "dev"
	cached: boolean
	sourcemap: boolean
	port: number
}

export interface ExportCommand {
	type: "export"
	cached: boolean
	sourcemap: boolean
}

export interface ServeCommand {
	type: "serve"
	port: number
}

export type AnyCommand = DevCommand | ExportCommand | ServeCommand

////////////////////////////////////////////////////////////////////////////////

// prettier-ignore
export interface Directories {
	wwwDir: string      // e.g. "www"
	srcPagesDir: string // e.g. "src/pages"
	cacheDir: string    // e.g. "__cache__"
	exportDir: string   // e.g. "__export__"
}

////////////////////////////////////////////////////////////////////////////////

// prettier-ignore
export interface FSPageInfo {
	type: "static" | "dynamic"
	src: string       // e.g. "src/pages/index.js"
	component: string // e.g. "PageIndex"
}

////////////////////////////////////////////////////////////////////////////////

export interface Props {
	[key: string]: unknown
}

export type DescriptProps = Props & { path: string }

export interface StaticPageModule {
	serverProps?(): Promise<Props>
	Head?: (props: DescriptProps) => JSX.Element
	default: (props: DescriptProps) => JSX.Element
}

export interface DynamicPageModule {
	serverPaths(): Promise<{ path: string; props: Props }[]>
	Head?: (props: DescriptProps) => JSX.Element
	default: (props: DescriptProps) => JSX.Element
}

export type AnyPageModule = StaticPageModule | DynamicPageModule

// prettier-ignore
export interface ServerRouteInfo {
	type: "static" | "dynamic"
	src: string       // e.g. "src/pages/index.js"
	dst: string       // e.g. "dst/index.html"
	path: string      // e.g. "/"
	component: string // e.g. "PageIndex"
}

export interface ServerRouteMeta {
	module: AnyPageModule
	route: ServerRouteInfo
	descriptProps: DescriptProps
}

export interface ServerRouter {
	[key: string]: ServerRouteMeta
}

////////////////////////////////////////////////////////////////////////////////

export interface Runtime<Cmd extends AnyCommand = AnyCommand> {
	cmd: Cmd
	dirs: Directories
	tmpl: string
	pages: FSPageInfo[]
	router: ServerRouter
	serverGuards(): Promise<void>
	purgeDirs(): Promise<void>
	resolveTemplate(): Promise<void>
	resolveFSPages(): Promise<void>
	resolveServerRouter(): Promise<void>
}
