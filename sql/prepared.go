package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"goeasy.dev/sql/middleware"
)

type namedStatement struct {
	chain middleware.Chain

	stmt *sqlx.NamedStmt
}

func (n namedStatement) Query(ctx context.Context, params interface{}) (rows *sqlx.Rows, err error) {
	n.chain.Exec(ctx, n.stmt.QueryString, []interface{}{params}, func(ctx context.Context) error {
		rows, err = n.stmt.QueryxContext(ctx, params)
		return err
	})

	return
}

func (n namedStatement) QueryRow(ctx context.Context, params interface{}) (row *sqlx.Row) {
	n.chain.Exec(ctx, n.stmt.QueryString, []interface{}{params}, func(ctx context.Context) error {
		row = n.stmt.QueryRowxContext(ctx, params)
		return nil
	})

	return
}

func (n namedStatement) Exec(ctx context.Context, params interface{}) (result sql.Result, err error) {
	n.chain.Exec(ctx, n.stmt.QueryString, []interface{}{params}, func(ctx context.Context) error {
		result, err = n.stmt.ExecContext(ctx, params)
		return err
	})

	return
}

func (n namedStatement) Select(ctx context.Context, dest interface{}, params interface{}) error {
	return n.chain.Exec(ctx, n.stmt.QueryString, []interface{}{params}, func(ctx context.Context) error {
		return n.stmt.SelectContext(ctx, dest, params)
	})
}

func (n namedStatement) Get(ctx context.Context, dest interface{}, params interface{}) error {
	return n.chain.Exec(ctx, n.stmt.QueryString, []interface{}{params}, func(ctx context.Context) error {
		return n.stmt.GetContext(ctx, dest, params)
	})
}

func (n namedStatement) Close() error {
	return n.stmt.Close()
}

type preparedStatment struct {
	chain middleware.Chain

	query string
	stmt  *sqlx.Stmt
}

func (p preparedStatment) Query(ctx context.Context, args ...interface{}) (rows *sqlx.Rows, err error) {
	p.chain.Exec(ctx, p.query, args, func(ctx context.Context) error {
		rows, err = p.stmt.QueryxContext(ctx, args...)
		return err
	})

	return
}

func (p preparedStatment) QueryRow(ctx context.Context, args ...interface{}) (row *sqlx.Row) {
	p.chain.Exec(ctx, p.query, args, func(ctx context.Context) error {
		row = p.stmt.QueryRowxContext(ctx, args...)
		return nil
	})

	return
}

func (p preparedStatment) Exec(ctx context.Context, args ...interface{}) (result sql.Result, err error) {
	p.chain.Exec(ctx, p.query, args, func(ctx context.Context) error {
		result, err = p.stmt.ExecContext(ctx, args...)
		return err
	})

	return
}

func (p preparedStatment) Select(ctx context.Context, dest interface{}, args ...interface{}) error {
	return p.chain.Exec(ctx, p.query, args, func(ctx context.Context) error {
		return p.stmt.SelectContext(ctx, dest, args...)
	})
}

func (p preparedStatment) Get(ctx context.Context, dest interface{}, args ...interface{}) error {
	return p.chain.Exec(ctx, p.query, args, func(ctx context.Context) error {
		return p.stmt.GetContext(ctx, dest, args...)
	})
}

func (p preparedStatment) Close() error {
	return p.stmt.Close()
}
