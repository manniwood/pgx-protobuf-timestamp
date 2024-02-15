[![Go Reference](https://pkg.go.dev/badge/github.com/manniwood/pgxtras.svg)](https://pkg.go.dev/github.com/manniwood/pgx-protobuf-timestamp)
![Build Status](https://github.com/manniwood/pgx-protobuf-timestamp/actions/workflows/ci.yml/badge.svg)

# pgx-protobuf-timestamp - Protobuf Timestamp Scanning for github.com/jackc/pgx

Are you using github.com/jackc/pgx/v5? Do you need to Scan/insert protobuf Timestamps?
Then this is the library for you!

## Setup

You will need to add this to your project of course:

```
go get github.com/manniwood/pgx-protobuf-timestamp
```

Then, if you are using a single connection to postgres, do this:

```
import (
	"github.com/jackc/pgx/v5"
	"github.com/manniwood/pgx-protobuf-timestamp/pgxpbts"
)

...

config, err := pgx.ParseConfig(dbURL)
if err != nil {
	return nil, err
}

conn, err := pgx.ConnectConfig(context.Background(), config)
if err != nil {
	return nil, err
}
pgxpbts.Register(conn.TypeMap())
```

or, if you are using pgxpool:

```
import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manniwood/pgx-protobuf-timestamp/pgxpbts"
)

...

config, err := pgxpool.ParseConfig(dbURL)
if err != nil {
	return nil, err
}
config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	pgxbpts.Register(conn.TypeMap())
	return nil
}
pool, err := pgxpool.NewWithConfig(context.Background(), config)
if err != nil {
	return nil, err
}
```

Now, you will be able to serialize/deserialize Go protobuf `*timestamppb.Timestamp` types
to/from Postgres `timestamp` and `timestamptz` types.

See the tests for examples.

