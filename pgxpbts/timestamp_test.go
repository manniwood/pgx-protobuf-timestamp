package pgxpbts_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/manniwood/pgx-protobuf-timestamp/pgxpbts"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const pgConnDSN string = "postgres://postgres:postgres@localhost:5432/postgres?application_name=BBB"

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
	want = time.Date(2024, time.February, 22, 12, 34, 56, 123456000, time.UTC)
	_, err = conn.Exec(ctx, `insert into t (id, ts) values (0, @ts)`, pgx.NamedArgs{"ts": timestamppb.New(want)})
	err = conn.QueryRow(ctx, `select ts from t where id = 0`).Scan(&got)
	if err != nil {
		t.Fatalf("Could not query db: %v", err)
	}
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

	return conn, nil
}
