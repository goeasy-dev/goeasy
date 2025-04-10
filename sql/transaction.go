package sql

import (
	"context"
	"database/sql"

	"goeasy.dev/errors"
	"goeasy.dev/sql/middleware"

	"github.com/jmoiron/sqlx"
)

type transaction struct {
	chain middleware.Chain

	tx *sqlx.Tx
}

func (c Connection) BeginTx(ctx context.Context, opts TxOptions) (Transaction, error) {
	tx, err := c.conn.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.IsolationLevel(opts.Isolation),
		ReadOnly:  opts.ReadOnly,
	})

	if err != nil {
		return nil, errors.Wrap(err, "unable to create db transaction")
	}

	return &transaction{
		chain: c.chain,
		tx:    tx,
	}, nil
}

func (t transaction) Commit() error {
	return t.tx.Commit()
}

func (t transaction) Rollback() error {
	return t.tx.Rollback()
}

// Query executes a SQL Query
func (t transaction) Query(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	return queryCtx(t.chain, t.tx, ctx, query, args...)
}

// Query Row executes a SQL Query
func (t transaction) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return queryRowCtx(t.chain, t.tx, ctx, query, args...)
}

// NamedQuery executes a SQL Query using the params to populate query parameters
func (t transaction) NamedQuery(ctx context.Context, query string, params interface{}) (rows *sqlx.Rows, err error) {
	return namedQueryCtx(t.chain, t.tx, ctx, query, params)
}

// Select executes a Select query, placing the results in dest
func (t transaction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return selectCtx(t.chain, t.tx, ctx, dest, query, args...)
}

// NamedSelect executes an SQL Query using the params to populate query parameters, placing the results in dest
func (t transaction) NamedSelect(ctx context.Context, dest interface{}, query string, params interface{}) error {
	return namedSelect(t.chain, t.tx, ctx, dest, query, params)
}

// Get executes a Select query, placing the first result in dest
func (t transaction) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return getCtx(t.chain, t.tx, ctx, dest, query, args...)
}

// NamedGet executes an SQL Query using the params to populate query parameters, placing the first result in dest
func (t transaction) NamedGet(ctx context.Context, dest interface{}, query string, params interface{}) error {
	return namedGet(t.chain, t.tx, ctx, dest, query, params)
}

func (t transaction) Exec(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	return execCtx(t.chain, t.tx, ctx, query, args...)
}

// NamedExec executes a sql query pulling arguments from the params
func (t transaction) NamedExec(ctx context.Context, query string, params interface{}) (result sql.Result, err error) {
	return namedExec(t.chain, t.tx, ctx, query, params)
}

func (t transaction) Prepare(ctx context.Context, query string) (Statement, error) {
	stmt, err := t.tx.PreparexContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare statement")
	}

	return &preparedStatment{
		chain: t.chain,
		stmt:  stmt,
	}, nil
}

func (t transaction) PrepareNamed(ctx context.Context, query string) (NamedStatement, error) {
	named, err := t.tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare named statement")
	}

	return &namedStatement{
		chain: t.chain,
		stmt:  named,
	}, nil
}
