package database

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	sqlx *sqlx.DB
}

type Tx struct {
	sqlx *sqlx.Tx
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
type emptyResult struct{}

func (rs emptyResult) LastInsertId() (int64, error) {
	return 0, nil
}
func (rs emptyResult) RowsAffected() (int64, error) {
	return 0, nil
}

// Need to add configuration for connection
func Connect(driver string, dataSource string) (*DB, error) {
	db, err := sqlx.Connect(driver, dataSource)
	if err != nil {
		return nil, err
	}

	// Need to add default connection configuration

	return &DB{sqlx: db}, nil
}

func (db *DB) bindNamed(query string, arg interface{}) (string, []interface{}, error) {
	if arg == nil {
		return query, []interface{}{}, nil
	}

	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return q, args, err
	}

	q, args, err = sqlx.In(q, args...)
	if err != nil {
		return q, args, err
	}
	q = db.sqlx.Rebind(q)

	return q, args, err
}

type testPing struct {
	Check string `db:"check"`
}

func (db *DB) Close() error {
	return db.sqlx.Close()
}

func (db *DB) Ping() error {
	dest := testPing{}
	return db.sqlx.GetContext(context.Background(), &dest, "SELECT 'connected' as check")
}

func (db *DB) Get(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := db.bindNamed(query, arg)
	if err != nil {
		return err
	}

	err = db.sqlx.GetContext(ctx, dest, q, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (db *DB) Select(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := db.bindNamed(query, arg)
	if err != nil {
		return err
	}

	err = db.sqlx.SelectContext(ctx, dest, q, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (db *DB) Exec(ctx context.Context, query string, arg interface{}) (Result, error) {
	q, args, err := db.bindNamed(query, arg)
	if err != nil {
		return emptyResult{}, err
	}

	return db.sqlx.ExecContext(ctx, q, args...)
}

func (db *DB) ExecWithArgs(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return db.sqlx.ExecContext(ctx, query, args...)
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.sqlx.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{sqlx: tx}, nil
}

func (tx *Tx) Get(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := tx.bindNamed(query, arg)
	if err != nil {
		return err
	}

	err = tx.sqlx.GetContext(ctx, dest, q, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (tx *Tx) Select(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := tx.bindNamed(query, arg)
	if err != nil {
		return err
	}

	err = tx.sqlx.SelectContext(ctx, dest, q, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (tx *Tx) Exec(ctx context.Context, query string, arg interface{}) (Result, error) {
	q, args, err := tx.bindNamed(query, arg)
	if err != nil {
		return emptyResult{}, err
	}

	return tx.sqlx.ExecContext(ctx, q, args...)
}

func (tx *Tx) ExecWithArgs(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return tx.sqlx.ExecContext(ctx, query, args...)
}

func (tx *Tx) Rollback() error {
	return tx.sqlx.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.sqlx.Commit()
}

func (tx *Tx) bindNamed(query string, arg interface{}) (string, []interface{}, error) {
	if arg == nil {
		return query, []interface{}{}, nil
	}

	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return q, args, err
	}

	q, args, err = sqlx.In(q, args...)
	if err != nil {
		return q, args, err
	}
	q = tx.sqlx.Rebind(q)

	return q, args, err
}
