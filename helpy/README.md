# Helpy - Toolkit of helpers

Those are written in Rust and are meant to be used as CLI tools.

Run the bun commands from the root of the project.

## Build

Compile all binaries

```bash
./build.sh
```

## Bunny Regions

See [API docs](https://docs.bunny.net/reference/regionpublic_index)

```bash
dotenvx run -- cargo run --bin bunny-regions
```

or

```bash
bun run bunny-regions
```

### Scaleway Runtimes

See [API docs](https://www.scaleway.com/en/developers/api/serverless-functions/#path-functions-list-function-runtimes)

```bash
dotenvx run -- cargo run --bin scw-runtimes
```

or

```bash
bun run scw-runtimes
```
