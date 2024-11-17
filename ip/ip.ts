import { Hono } from 'https://esm.sh/jsr/@hono/hono@4.6.10'
import { cors } from 'https://esm.sh/jsr/@hono/hono@4.6.10/cors'
import { etag } from 'https://esm.sh/jsr/@hono/hono@4.6.10/etag'
import { prettyJSON } from 'https://esm.sh/jsr/@hono/hono@4.6.10/pretty-json'

// create hono app server
const app = new Hono()
app.use(prettyJSON({ space: 4 }))
app.use('*', etag())
app.use('*', cors())

// check if ip is ipv4 or ipv6
const isIPv4 = (ip: string) => ip.indexOf(':') === -1
const isIPv6 = (ip: string) => ip.indexOf(':') !== -1

// disable caching
app.use('*', (c, next) => {
	c.header('cache-control', 'private; max-age=0; must-revalidate')
	return next()
})

// catch all requests
app.get('/*', (c) =>
	c.json({
		isIpv4: isIPv4(c.req.header('x-real-ip') ?? ''),
		isIpv6: isIPv6(c.req.header('x-real-ip') ?? ''),
		ip: c.req.header('x-real-ip'),
		method: c.req.method,
		url: c.req.url,
		header: c.req.header(),
		query: c.req.query(),
	})
)

// provide export for Bunny runtime
// @ts-expect-error: <Bunny is available in the EdgeScript environment>
// biome-ignore lint/correctness/noUndeclaredVariables: <Bunny is available in the EdgeScript environment>
if (typeof Bunny !== 'undefined') Bunny.v1.serve((req: Request): Response | Promise<Response> => app.fetch(req))

// export for other runtimes
export default app
