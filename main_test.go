package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	_ "github.com/joho/godotenv/autoload"
)

const MaxRows = 10_000

var (
	conn *pgx.Conn
	ctx  context.Context
	err  error
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func BenchmarkInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanTable()
		for c := 0; c < MaxRows; c++ {
			query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("test").
				Columns("id", "name", "meta", "status", "created_at", "updated_at")
			query = query.Values(c+1, "name"+strconv.Itoa(c), "", "NEW", time.Now(), time.Now())
			q, args, _ := query.ToSql()
			_, err = conn.Exec(ctx, q, args...)
			if err != nil {
				fatal("cannot insert to table: %v\n", err)
			}
		}
	}
}

func BenchmarkTransactionInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanTable()
		for c := 0; c < MaxRows; c++ {
			query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("test").
				Columns("id", "name", "meta", "status", "created_at", "updated_at")
			query = query.Values(c+1, "name"+strconv.Itoa(c), "", "NEW", time.Now(), time.Now())
			q, args, _ := query.ToSql()
			err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
				_, err = tx.Exec(ctx, q, args...)
				return err
			})
			if err != nil {
				fatal("cannot insert to table: %v\n", err)
			}
		}
	}
}

func BenchmarkBulkInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanTable()
		query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("test").
			Columns("id", "name", "meta", "status", "created_at", "updated_at")
		for c := 0; c < MaxRows; c++ {
			query = query.Values(c+1, "name"+strconv.Itoa(c), "", "NEW", time.Now(), time.Now())
		}
		q, args, _ := query.ToSql()
		_, err = conn.Exec(ctx, q, args...)
		if err != nil {
			fatal("cannot insert to table: %v\n", err)
		}
	}
}

func BenchmarkTransactionBulkInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanTable()
		query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("test").
			Columns("id", "name", "meta", "status", "created_at", "updated_at")
		for c := 0; c < MaxRows; c++ {
			query = query.Values(c+1, "name"+strconv.Itoa(c), "", "NEW", time.Now(), time.Now())
		}
		q, args, _ := query.ToSql()
		conn.BeginFunc(ctx, func(tx pgx.Tx) error {
			_, err = conn.Exec(ctx, q, args...)
			return err
		})
		if err != nil {
			fatal("cannot insert to table: %v\n", err)
		}
	}
}

func BenchmarkCopyFromInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanTable()
		var rows [][]any
		for c := 0; c < MaxRows; c++ {
			rows = append(rows, []any{c + 1, "name" + strconv.Itoa(c), "", "NEW", time.Now(), time.Now()})
		}
		_, err = conn.CopyFrom(
			ctx,
			pgx.Identifier{"test"},
			[]string{"id", "name", "meta", "status", "created_at", "updated_at"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			fatal("cannot insert to table: %v\n", err)
		}
	}
}

func BenchmarkTransactionCopyFromInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanTable()
		var rows [][]any
		for c := 0; c < MaxRows; c++ {
			rows = append(rows, []any{c + 1, "name" + strconv.Itoa(c), "", "NEW", time.Now(), time.Now()})
		}
		err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
			_, err = conn.CopyFrom(
				ctx,
				pgx.Identifier{"test"},
				[]string{"id", "name", "meta", "status", "created_at", "updated_at"},
				pgx.CopyFromRows(rows),
			)
			return err
		})
		if err != nil {
			fatal("cannot insert to table: %v\n", err)
		}
	}
}

func setup() {
	// setup database
	ctx = context.Background()
	cfg, err := pgx.ParseConfig(os.Getenv("PGX_TEST_DATABASE"))
	if err != nil {
		fatal("unable to parse config: %v\n", err)
	}
	conn, err = pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		fatal("unable to connect to database: %v\n", err)
	}

	_, err = conn.Exec(ctx, `drop table IF EXISTS test; 
	create table test(
		id bigint NOT NULL,
		name text,
		meta text,
		status text,
		created_at timestamp,
		updated_at timestamp
	);
	create index on test(id);`)
	if err != nil {
		fatal("cannot create table: %v\n", err)
	}
}

func teardown() {
	// remove table
	_, err = conn.Exec(ctx, `drop table test;`)
	if err != nil {
		fatal("cannot drop table: %v\n", err)
	}
}

func fatal(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func cleanTable() {
	_, err = conn.Exec(ctx, `TRUNCATE TABLE test;`)
	if err != nil {
		fatal("cannot truncate table: %v\n", err)
	}
	_, err = conn.Exec(ctx, `vacuum test;`)
	if err != nil {
		fatal("cannot vacuum table: %v\n", err)
	}
}
