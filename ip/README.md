# frytg / ip

This simple script returns the IP address and other parameters of a request. It is designed as a test for using [Deno](https://deno.com) + [Hono](https://hono.dev) via [esm.sh](https://esm.sh).

## Usage

```bash
deno serve --port 8001 ip/ip.ts
```

Then open [`localhost:8001`](http://localhost:8001/) in your browser. You can also add [`?pretty`](http://localhost:8001/?pretty) to get a pretty-printed JSON response.

## Deployment

You can throw this on any FaaS that supports either Deno or `esm.sh`. It is currently being used on [Bunny](https://bunny.net?ref=qb6g5ox5lv) EdgeScript.
