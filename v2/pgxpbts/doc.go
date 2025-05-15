// Package pgxpbts is used to help pgx scan postgres timestamp
// and timestamptz types into Google Protobuf type *timestamppb.Timestamp.
//
// For a single pgx connection, *pgx.Conn, register both type translators
// like so:
//
//	pgxpbts.Register(conn.TypeMap())
//	pgxpbts.RegisterTZ(conn.TypeMap())
//
// For *pgxpool.Pool, register both type translators like so:
//
//	config, err := pgxpool.ParseConfig(dbURL)
//	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
//		pgxpbts.Register(conn.TypeMap())
//		pgxpbts.RegisterTZ(conn.TypeMap())
//		return nil
//	}
//	pool, err := pgxpool.NewWithConfig(context.Background(), config)
//
// Now, you will be able to serialize/deserialize Go protobuf *timestamppb.Timestamp types
// to/from Postgres timestamp and timestamptz types.
//
// See the tests for examples.
//
// # Inspiration
//
// See https://github.com/jackc/pgx-gofrs-uuid and
// https://github.com/jackc/pgx-shopspring-decimal for more examples of
// how to write type serializers/deserializers for pgx.
package pgxpbts
