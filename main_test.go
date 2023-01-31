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

const (
	MaxRows       = 10_000
	MaxRowsUpdate = 3_000
)

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
	for i := 0; i < b.N; i++ {
		cleanTable()
		query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("test").
			Columns("id", "name", "meta", "status", "created_at", "updated_at")
		for c := 0; c < MaxRows; c++ {
			query = query.Values(c+1, "name"+strconv.Itoa(c), "", "NEW", time.Now(), time.Now())
		}
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

func BenchmarkCopyFromInsert(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		cleanTable()
		var rows [][]any
		for c := 0; c < MaxRows; c++ {
			rows = append(rows, []any{c + 1, "name" + strconv.Itoa(c), "", "NEW", time.Now(), time.Now()})
		}
		err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
			_, err = tx.CopyFrom(
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

// TODO: bulk update for 30_000_000 elements
func BenchmarkUpdate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cleanTable()
		fillTable()
		var sent *time.Time
		t := time.Now()
		sent = &t
		var ids []uint64
		for c := 0; c < MaxRowsUpdate; c++ {
			ids = append(ids, uint64(c))
		}
		qr, args, _ := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update("test").
			Set("status", "SENT").
			Set("created_at", sent).
			Set("updated_at", time.Now()).
			Where(sq.Eq{"id": ids}).ToSql()
		if _, err = conn.Exec(ctx, qr, args...); err != nil {
			fatal("cannot insert to table: %v\n", err)
		}
	}
}
func BenchmarkUpdateWithTemporaryTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cleanTable()
		fillTable()
		var rows [][]any
		for c := 0; c < MaxRowsUpdate; c++ {
			var sent *time.Time
			t := time.Now()
			sent = &t
			rows = append(rows, []any{c + 1, "name" + strconv.Itoa(c), "", "SENT", sent, time.Now()})
		}
		err = conn.BeginFunc(ctx, func(tx pgx.Tx) error {
			if _, err = tx.Exec(ctx,
				`create temporary table tmp(
						id bigint,
						name text,
						meta text,
						status text,
						created_at timestamp,
						updated_at timestamp)`,
			); err != nil {
				return err
			}

			if _, err = tx.CopyFrom(
				ctx,
				pgx.Identifier{"tmp"},
				[]string{"id", "name", "meta", "status", "created_at", "updated_at"},
				pgx.CopyFromRows(rows),
			); err != nil {
				return err
			}

			if _, err = tx.Exec(ctx, `CREATE INDEX ON tmp(id)`); err != nil {
				return err
			}
			// update main table
			if _, err = tx.Exec(ctx,
				`update test 
					SET 
						status=t.status, 
						updated_at=t.updated_at, 
						created_at=t.created_at 
					FROM tmp t 
					WHERE 
						t.id=test.id
					AND 
						t.name=test.name`,
			); err != nil {
				return err
			}
			// drop temporary table
			if _, err = tx.Exec(ctx, `drop table tmp;`); err != nil {
				return err
			}
			return nil
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

func fillTable() {
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
