package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"goeasy.dev/sql/middleware"
)

type Conn interface {
	Querier
	Executor
	Preparer

	BeginTx(ctx context.Context, opts TxOptions) (Transaction, error)
}

type Transaction interface {
	Querier
	Executor
	Preparer

	Commit() error
	Rollback() error
}

type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedQuery(ctx context.Context, query string, params interface{}) (rows *sqlx.Rows, err error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedSelect(ctx context.Context, dest interface{}, query string, params interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedGet(ctx context.Context, dest interface{}, query string, params interface{}) error
}

type Executor interface {
	Exec(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error)
	NamedExec(ctx context.Context, query string, params interface{}) (result sql.Result, err error)
}

type Preparer interface {
	Prepare(ctx context.Context, query string) (Statement, error)
	PrepareNamed(ctx context.Context, query string) (NamedStatement, error)
}

type Statement interface {
	Query(ctx context.Context, args ...interface{}) (rows *sqlx.Rows, err error)
	QueryRow(ctx context.Context, args ...interface{}) *sqlx.Row
	Exec(ctx context.Context, args ...interface{}) (result sql.Result, err error)
	Get(ctx context.Context, dest interface{}, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, args ...interface{}) error
	Close() error
}

type NamedStatement interface {
	Query(ctx context.Context, params interface{}) (rows *sqlx.Rows, err error)
	QueryRow(ctx context.Context, params interface{}) *sqlx.Row
	Exec(ctx context.Context, params interface{}) (result sql.Result, err error)
	Get(ctx context.Context, dest interface{}, params interface{}) error
	Select(ctx context.Context, dest interface{}, params interface{}) error
	Close() error
}

func queryCtx(m middleware.Chain, db sqlx.QueryerContext, ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	err = m.Exec(ctx, query, args, func(ctx context.Context) error {
		rows, err = db.QueryxContext(ctx, query, args...)
		return err
	})

	return
}

func queryRowCtx(m middleware.Chain, db sqlx.QueryerContext, ctx context.Context, query string, args ...interface{}) (row *sqlx.Row) {
	m.Exec(ctx, query, args, func(ctx context.Context) error {
		row = db.QueryRowxContext(ctx, query, args...)
		return nil
	})

	return
}

func namedQueryCtx(m middleware.Chain, db sqlx.ExtContext, ctx context.Context, query string, params interface{}) (rows *sqlx.Rows, err error) {
	err = m.Exec(ctx, query, []interface{}{params}, func(ctx context.Context) error {
		rows, err = sqlx.NamedQueryContext(ctx, db, query, params)
		return err
	})

	return
}

func selectCtx(m middleware.Chain, db sqlx.QueryerContext, ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.Exec(ctx, query, args, func(ctx context.Context) error {
		return sqlx.SelectContext(ctx, db, dest, query, args...)
	})
}

func getCtx(m middleware.Chain, db sqlx.QueryerContext, ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.Exec(ctx, query, args, func(ctx context.Context) error {
		return sqlx.GetContext(ctx, db, dest, query, args...)
	})
}

func namedGet(m middleware.Chain, db sqlx.ExtContext, ctx context.Context, dest interface{}, query string, params interface{}) error {
	return m.Exec(ctx, query, []interface{}{params}, func(ctx context.Context) error {
		rows, err := sqlx.NamedQueryContext(ctx, db, query, params)
		if err != nil {
			return err
		}
		defer rows.Close()

		if rows.Next() {
			err = rows.StructScan(dest)
			if err != nil {
				return err
			}
		}

		return rows.Err()
	})
}

func namedSelect(m middleware.Chain, db sqlx.ExtContext, ctx context.Context, dest interface{}, query string, params interface{}) error {
	return m.Exec(ctx, query, []interface{}{params}, func(ctx context.Context) error {
		rows, err := sqlx.NamedQueryContext(ctx, db, query, params)
		if err != nil {
			return err
		}
		defer rows.Close()

		return sqlx.StructScan(rows, dest)
	})
}

func execCtx(m middleware.Chain, db sqlx.ExecerContext, ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	err = m.Exec(ctx, query, args, func(ctx context.Context) error {
		result, err = db.ExecContext(ctx, query, args...)
		return err
	})

	return
}

func namedExec(m middleware.Chain, db sqlx.ExtContext, ctx context.Context, query string, params interface{}) (result sql.Result, err error) {
	err = m.Exec(ctx, query, []interface{}{params}, func(ctx context.Context) error {
		result, err = sqlx.NamedExecContext(ctx, db, query, params)
		return err
	})

	return
}
