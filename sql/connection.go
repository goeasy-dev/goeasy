package sql

import (
	"context"
	"database/sql"
	"time"

	"goeasy.dev/errors"
	"goeasy.dev/observability/log"
	"goeasy.dev/sql/middleware"
	"goeasy.dev/status"
	"goeasy.dev/status/statustype"
	"goeasy.dev/util"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	watchInterval       = time.Second * 3
	failedWatchInterval = time.Millisecond * 500
)

var _ Conn = Connection{}

// Config contains the information for creating a Connection
type Config struct {
	Driver           string
	LogQueries       bool
	ConnectionString ConnectionString
}

// connection constitutes a connection to a database
type Connection struct {
	chain middleware.Chain

	conn       *sqlx.DB
	config     Config
	closedChan chan struct{}
}

// NewConnection creates a new database connection
func NewConnection(ctx context.Context, config Config) (Connection, error) {
	connection := Connection{
		chain: middleware.Chain{
			middleware.LogMiddleware,
			middleware.MetricsMiddleware,
			middleware.TraceMiddleware,
		},
		conn:       &sqlx.DB{},
		closedChan: make(chan struct{}),
		config:     config,
	}

	dsn, err := config.ConnectionString.Build(ctx)
	if err != nil {
		return connection, errors.Wrap(err, "unable to build connection string")
	}

	conn, err := sqlx.ConnectContext(ctx, config.Driver, dsn)
	if err != nil {
		return connection, errors.Wrap(err, "unable to create db connection")
	}

	connection.conn = conn

	if err = connection.Ping(); err != nil {
		return connection, errors.Wrap(err, "unable to connect to db")
	}

	go func() {
		<-ctx.Done()
		connection.Close()
	}()

	go watch(ctx, connection)

	return connection, nil
}

// Reestablish the connection with the database
func (c Connection) Reestablish(ctx context.Context) error {
	dsn, err := c.config.ConnectionString.Build(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to build connection string")
	}

	conn, err := sqlx.ConnectContext(ctx, c.config.Driver, dsn)
	if err != nil {
		return errors.Wrap(err, "unable to create db connection")
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "unable to close existing connection"))
		}
	}

	*c.conn = *conn
	return nil
}

func (c Connection) DB() *sql.DB {
	return c.conn.DB
}

func (c Connection) DBx() *sqlx.DB {
	return c.conn
}

func (c Connection) Ping() error {
	return c.conn.Ping()
}

func (c Connection) Close() error {
	c.closedChan <- struct{}{}

	err := c.conn.Close()
	return err
}

// Query executes a SQL Query
func (c Connection) Query(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	return queryCtx(c.chain, c.conn, ctx, query, args...)
}

// Query executes a SQL Query
func (c Connection) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return queryRowCtx(c.chain, c.conn, ctx, query, args...)
}

// NamedQuery executes a SQL Query using the params to populate query parameters
func (c Connection) NamedQuery(ctx context.Context, query string, params interface{}) (rows *sqlx.Rows, err error) {
	return namedQueryCtx(c.chain, c.conn, ctx, query, params)
}

// Select executes a Select query, placing the results in dest
func (c Connection) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return selectCtx(c.chain, c.conn, ctx, dest, query, args...)
}

// NamedSelect executes an SQL Query using the params to populate query parameters, placing the results in dest
func (c Connection) NamedSelect(ctx context.Context, dest interface{}, query string, params interface{}) error {
	return namedSelect(c.chain, c.conn, ctx, dest, query, params)
}

// Get executes a Select query, placing the first result in dest
func (c Connection) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return getCtx(c.chain, c.conn, ctx, dest, query, args...)
}

// NamedGet executes an SQL Query using the params to populate query parameters, placing the first result in dest
func (c Connection) NamedGet(ctx context.Context, dest interface{}, query string, params interface{}) error {
	return namedGet(c.chain, c.conn, ctx, dest, query, params)
}

// Exec executes the sql query
func (c Connection) Exec(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	return execCtx(c.chain, c.conn, ctx, query, args...)
}

// NamedExec executes a sql query pulling arguments from the params
func (c Connection) NamedExec(ctx context.Context, query string, params interface{}) (result sql.Result, err error) {
	return namedExec(c.chain, c.conn, ctx, query, params)
}

func (c Connection) Prepare(ctx context.Context, query string) (Statement, error) {
	stmt, err := c.conn.PreparexContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare named statement")
	}

	return &preparedStatment{
		chain: c.chain,
		stmt:  stmt,
	}, nil
}

func (c Connection) PrepareNamed(ctx context.Context, query string) (NamedStatement, error) {
	named, err := c.conn.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare named statement")
	}

	return &namedStatement{
		chain: c.chain,
		stmt:  named,
	}, nil
}

func watch(ctx context.Context, c Connection) {
	check := status.SimpleCheck(statustype.Liveness & statustype.Readiness)
	check = util.Ptr(true)

	isFailed := false
	for {
		select {
		case <-c.closedChan:
			return
		default:
			log.Trace("sql: validating db connection")
		}

		isFailed = !liveLoopHandler(ctx, c, check, isFailed)

		if !isFailed {
			time.Sleep(watchInterval)
		} else {
			time.Sleep(failedWatchInterval)
		}
	}
}

func liveLoopHandler(ctx context.Context, c Connection, check *bool, isFailed bool) bool {
	if !isFailed {
		if checkIfLive(ctx, c) {
			check = util.Ptr(true)
			return true
		}

		check = util.Ptr(false)
	}

	log.Trace("sql: attempting db reestablish")
	check = util.Ptr(false)
	if reestablish(ctx, c) {
		log.Trace("sql: connection reestablished")
		return true
	}

	return false
}

func checkIfLive(ctx context.Context, c Connection) bool {
	var val string
	if err := c.Get(ctx, &val, "SELECT current_user"); err != nil {
		log.Warn(errors.Wrap(err, "ping failed"))
		return false
	}

	log.Trace("sql: db ping successful")
	return true
}

func reestablish(ctx context.Context, c Connection) bool {
	err := c.Reestablish(ctx)
	if err != nil {
		log.Error(errors.Wrap(err, "unable to reestablish connection"))
		return false
	}

	log.Trace("sql: connection reestablished")
	return true
}
