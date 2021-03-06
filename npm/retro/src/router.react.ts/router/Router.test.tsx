import React from "react"
import renderer from "react-test-renderer"
import { BrowserRouter } from "./BrowserRouter"
import { Route } from "./Route"
import { Router } from "./Router"

beforeAll(() => {
	mockLocationPathname()
})

function RoutedApp() {
	return (
		<BrowserRouter>
			<Router>
				<Route path="/">
					<h1>
						Hello, <code>/</code>!
					</h1>
				</Route>
				<Route path="/page-a">
					<h1>
						Hello, <code>/page-a</code>!
					</h1>
				</Route>
				<Route path="/page-b">
					<h1>
						Hello, <code>/page-b</code>!
					</h1>
				</Route>
				<Route path="/404">
					<h1>
						Hello, <code>/404</code>!
					</h1>
				</Route>
			</Router>
		</BrowserRouter>
	)
}

test("should route to /", () => {
	window.location.pathname = "/"
	const tree = renderer.create(<RoutedApp />).toJSON()
	expect(tree).toMatchInlineSnapshot(`
    <h1>
      Hello,\u0020
      <code>
        /
      </code>
      !
    </h1>
  `)
})

test("should route to /page-a", () => {
	window.location.pathname = "/page-a"
	const tree = renderer.create(<RoutedApp />).toJSON()
	expect(tree).toMatchInlineSnapshot(`
    <h1>
      Hello,\u0020
      <code>
        /page-a
      </code>
      !
    </h1>
  `)
})

test("should route to /page-b", () => {
	window.location.pathname = "/page-b"
	const tree = renderer.create(<RoutedApp />).toJSON()
	expect(tree).toMatchInlineSnapshot(`
    <h1>
      Hello,\u0020
      <code>
        /page-b
      </code>
      !
    </h1>
  `)
})

test("should route to /404", () => {
	window.location.pathname = "/404"
	const tree = renderer.create(<RoutedApp />).toJSON()
	expect(tree).toMatchInlineSnapshot(`
    <h1>
      Hello,\u0020
      <code>
        /404
      </code>
      !
    </h1>
  `)
})
