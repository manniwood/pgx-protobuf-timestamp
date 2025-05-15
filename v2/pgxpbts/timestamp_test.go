package pgxpbts_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/manniwood/pgx-protobuf-timestamp/v2/pgxpbts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const pgConnDSN string = "postgres://postgres:postgres@localhost:5432/postgres?application_name=BBB"

// timestamppb.Timestamp has no concept of infinity;
// so when we read 'infinity' or '-infinity' timestamps from Postgres,
// this is the "magic date" we will get.
const infinityMagicDateStr string = "0001-01-01 00:00:00.000000000"

func TestTimestamp(t *testing.T) {
	ctx := context.Background()
	conn, err := newConn(ctx, pgConnDSN)
	if err != nil {
		t.Fatalf("Could not connect to db: %v", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `create temporary table t (id int constraint t_is_singleton check (id = 0) not null, ts timestamp not null)`)
	if err != nil {
		t.Fatalf("Could not create temporary table: %v", err)
	}

	var got *timestamppb.Timestamp
	var want time.Time
	wantStr := "2024-01-22 12:34:56.123456000"
	want, err = time.Parse("2006-01-02 03:04:05.999999999", wantStr)
	if err != nil {
		t.Fatalf("Could not parse time %s: %v", wantStr, err)
	}
	_, err = conn.Exec(ctx, `insert into t (id, ts) values (0, @ts)`, pgx.NamedArgs{"ts": timestamppb.New(want)})
	if err != nil {
		t.Fatalf("Could not insert: %v", err)
	}
	err = conn.QueryRow(ctx, `select ts from t where id = 0`).Scan(&got)
	if err != nil {
		t.Fatalf("Could not query db: %v", err)
	}
	if !got.AsTime().Equal(want) {
		t.Errorf("Got %s; want %s", got.AsTime(), want)
	}
}

func TestNilTimestamp(t *testing.T) {
	ctx := context.Background()
	conn, err := newConn(ctx, pgConnDSN)
	if err != nil {
		t.Fatalf("Could not connect to db: %v", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `create temporary table t (id int constraint t_is_singleton check (id = 0) not null, ts timestamp)`)
	if err != nil {
		t.Fatalf("Could not create temporary table: %v", err)
	}

	var input *timestamppb.Timestamp
	var got *timestamppb.Timestamp
	// _, err = conn.Exec(ctx, `insert into t (id, ts) values (0, @ts)`, pgx.NamedArgs{"ts": nil})
	_, err = conn.Exec(ctx, `insert into t (id, ts) values (0, @ts)`, pgx.NamedArgs{"ts": input})
	if err != nil {
		t.Fatalf("Could not insert: %v", err)
	}
	err = conn.QueryRow(ctx, `select ts from t where id = 0`).Scan(&got)
	if err != nil {
		t.Fatalf("Could not query db: %v", err)
	}
	if got != nil {
		t.Errorf("Got %s; want nil", got.AsTime())
	}
}

func TestInfiniteTimestamp(t *testing.T) {
	ctx := context.Background()
	conn, err := newConn(ctx, pgConnDSN)
	if err != nil {
		t.Fatalf("Could not connect to db: %v", err)
	}
	defer conn.Close(ctx)

	var got *timestamppb.Timestamp
	err = conn.QueryRow(ctx, `select timestamp 'infinity'`).Scan(&got)
	if err != nil {
		t.Fatalf("Could not query db: %v", err)
	}
	var want time.Time
	want, err = time.Parse("2006-01-02 03:04:05.999999999", infinityMagicDateStr)
	if !got.AsTime().Equal(want) {
		t.Errorf("Got %s; want %s", got.AsTime(), want)
	}
}

func TestNegativeInfiniteTimestamp(t *testing.T) {
	ctx := context.Background()
	conn, err := newConn(ctx, pgConnDSN)
	if err != nil {
		t.Fatalf("Could not connect to db: %v", err)
	}
	defer conn.Close(ctx)

	var got *timestamppb.Timestamp
	err = conn.QueryRow(ctx, `select timestamp '-infinity'`).Scan(&got)
	if err != nil {
		t.Fatalf("Could not query db: %v", err)
	}
	var want time.Time
	want, err = time.Parse("2006-01-02 03:04:05.999999999", infinityMagicDateStr)
	if !got.AsTime().Equal(want) {
		t.Errorf("Got %s; want %s", got.AsTime(), want)
	}
}

func newConn(ctx context.Context, dbURL string) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	pgxpbts.Register(conn.TypeMap())
	pgxpbts.RegisterTZ(conn.TypeMap())

	return conn, nil
}
