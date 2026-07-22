# golang-scenarios

A small service for authoring and running scripted evaluation scenarios. A scenario spec is a JSON file describing a
set of evals — free-form expectations plus structured assertions — and the `scenarios` binary loads, normalizes, and
reports on them. A minimal web front end (sign-in via an OAuth provider, with a small SPA status page) is included for
browsing runs.

## Usage

Inspect a spec file:

```sh
go run ./cmd/scenarios spec examples/greeting.json
```

Run the web UI:

```sh
go run ./cmd/scenarios serve :8080
```

The server exposes `/login`, `/callback`, and `/logout`. Sign-in state is carried in the `scenarios-auth` session
cookie; the SPA reads the `scenarios-auth-provider` cookie to render the sign-in surface.

## Spec format

```json
{
  "evals": [
    {
      "name": "greeting",
      "expectations": ["is polite", "mentions the user by name"],
      "assertions": [{ "kind": "contains", "value": "hello" }]
    }
  ]
}
```

On load, each eval's `expectations` are expanded into leading `llm` assertions and merged ahead of the authored
`assertions`, so downstream execution handles a single list.

## Development

```sh
go build ./...
go vet ./...
go test -race ./...
```

CI runs the same gates plus golangci-lint; CodeQL analyzes the module and the workflows on every push and pull
request.
