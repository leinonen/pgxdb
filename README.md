# pgxdb

A [glisp](https://github.com/leinonen/golisp-language) module wrapping the [pgx v5](https://github.com/jackc/pgx) PostgreSQL driver.

## Installation

```
glisp get github.com/leinonen/pgxdb@v0.1.0
```

This also wires `github.com/jackc/pgx/v5` into your project's `go.mod` automatically via `go-require`.

## Usage

```clojure
(ns main
  (:import [fmt])
  (:require [github.com/leinonen/pgxdb]))

(defn -main []
  (if-err [conn err] (pgxdb/connect "postgres://user:pass@localhost:5432/mydb")
    (fmt/println "connect failed:" (.Error err))
    (do
      ;; INSERT / UPDATE / DELETE — returns (rows-affected, error)
      (if-err [n err] (pgxdb/exec conn
                        "INSERT INTO users (name, email) VALUES ($1, $2)"
                        ["Alice" "alice@example.com"])
        (fmt/println "insert failed:" (.Error err))
        (fmt/println "inserted" n "rows"))

      ;; SELECT — returns a vector of row maps, each map keyed by column name
      (if-err [rows err] (pgxdb/query conn
                           "SELECT * FROM users WHERE name = $1"
                           ["Alice"])
        (fmt/println "query failed:" (.Error err))
        (fmt/println "result:" rows))

      (pgxdb/close conn))))
```

The connection URL is read from the `DATABASE_URL` environment variable in practice:

```clojure
(def db-url (os/env "DATABASE_URL" "postgres://localhost/mydb"))
```

## API

| Function | Signature | Description |
|---|---|---|
| `pgxdb/connect` | `url → (conn, error)` | Open a connection. `conn` is opaque — pass it to the other functions. |
| `pgxdb/close` | `conn → error` | Close the connection. |
| `pgxdb/query` | `conn sql args → (rows, error)` | Run a SELECT. Returns a vector of maps `{column → value}`. Pass `nil` for `args` when no parameters. |
| `pgxdb/exec` | `conn sql args → (rows-affected, error)` | Run INSERT / UPDATE / DELETE. Pass `nil` for `args` when no parameters. |

Use `$1`, `$2`, … placeholders for query parameters — never interpolate values directly into SQL strings.

## Running the example

The `example/` directory contains a full demo app with a Docker Compose postgres instance.

```bash
cd example

# start postgres
docker compose up -d

# build and run (glisp build handles dependency wiring automatically)
make run
```

Expected output (first run):

```
table ready
inserted 2 rows
all users:
   1 Alice alice@example.com
   2 Bob bob@example.com
alice only:
   1 Alice alice@example.com
done
```

Subsequent runs show `inserted 0 rows` — the `ON CONFLICT DO NOTHING` clause makes inserts idempotent.

## How it works

glisp modules can declare Go package dependencies in `glisp.mod` using `go-require`:

```
module github.com/leinonen/pgxdb

go-require (
  github.com/jackc/pgx/v5 v5.7.2
)
```

When a project runs `glisp get github.com/leinonen/pgxdb`, the toolchain:

1. Downloads and transpiles the module's `.glsp` files
2. Writes `go-require` entries into the module's own `go.mod`
3. Propagates them into the project's `go.mod` — so `go build` can find pgx

The module itself uses a **bridge pattern**: a hand-written `bridge.go` handles the
variadic pgx API and type assertions; `db.glsp` exposes the clean glisp-facing API on top.
