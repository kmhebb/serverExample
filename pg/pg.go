package pg

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"

	cloud "github.com/kmhebb/serverExample"
	"github.com/kmhebb/serverExample/log"
)

type Tx interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Commit(ctx context.Context) error
}

type tx struct {
	tx pgx.Tx
}

func (tx tx) Commit(ctx context.Context) error {
	return tx.tx.Commit(ctx)
}

func (tx tx) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := tx.tx.Exec(ctx, query, args...)
	return err
}

func (tx tx) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := tx.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (tx tx) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return tx.tx.QueryRow(ctx, query, args)
}

func NewDatabase(ctx context.Context, url string) (*Database, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Printf("url: %s", url)
		return nil, fmt.Errorf("pg/NewDatabase: %w", err)
	}
	db := &Database{
		conn: conn,
	}
	return db, nil
}

type Database struct {
	conn *pgx.Conn
}

func (db Database) execFile(ctx context.Context, path string) error {
	log.Debug("running migration", log.Fields{"path": path})

	migration, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("pg/Database.ExecFile: %w", err)
	}

	if _, err := db.conn.Exec(ctx, string(migration)); err != nil {
		return fmt.Errorf("pg/Database.ExecFile: %w", err)
	}

	return nil
}

func (db *Database) Close(ctx context.Context) error {
	if err := db.conn.Close(ctx); err != nil {
		return fmt.Errorf("pg/Database.Close: %w", err)
	}
	return nil
}

func (db *Database) Ping(ctx context.Context) error {
	if db.conn.IsClosed() {
		log.Info("connection closed. attempting to reconnect.", nil)
		conn, err := pgx.Connect(ctx, db.conn.Config().ConnString())
		if err != nil {
			return fmt.Errorf("pg/Database.Ping: %w", err)
		}
		db.conn = conn
		return nil
	}

	if err := db.conn.Ping(ctx); err != nil {
		return fmt.Errorf("pg/Database.Ping: %w", err)
	}

	return nil
}

func (db *Database) RunInTransaction(ctx cloud.Context, f func(cloud.Context, Tx) error) error {
	err := db.conn.BeginFunc(ctx.Ctx, func(pgxtx pgx.Tx) error {
		tx := tx{tx: pgxtx}
		return f(ctx, tx)
	})
	if err != nil {
		return fmt.Errorf("pg/Database.RunIntransaction: %w", err)
	}
	return nil
}
