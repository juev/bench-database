package main

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var (
	db  *sql.DB
	err error
)

func fatal(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func BenchmarkInsert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for c := 0; c < 10_000; c++ {
			_, _ = db.Exec(`insert into test (id, name, meta, status, created_at, updated_at) VALUES
						(c+1, "name"+strconv.Itoa(c),"", "NEW", time.NOW(), time.NOW());`)
		}
		_, err = db.Exec(`TRUNCATE TABLE test;`)
		if err != nil {
			fatal("cannot truncate table: %v", err)
		}
		_, err = db.Exec(`vacuum test;`)
		if err != nil {
			fatal("cannot vacuum table: %v", err)
		}
	}
}

func setup() {
	// setup database
	db, err = sql.Open("pgx", os.Getenv("PGX_TEST_DATABASE"))
	if err != nil {
		fatal("unable to connect to database: %v\n", err)
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	if err := goose.Up(db, "migrations"); err != nil {
		fatal("error during migrations: %v\n", err)
	}
}

func teardown() {
	// remove table
	if err := goose.Down(db, "migrations"); err != nil {
		fatal("error during migrations: %v\n", err)
	}
	db.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
