package sql_test

import (
	"context"
	dbsql "database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"goeasy.dev/sql"
)

type testRecord struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type ConnectionTestSuite struct {
	suite.Suite
	conn sql.Connection
	ctx  context.Context
}

func (s *ConnectionTestSuite) SetupSuite() {
	s.ctx = context.Background()
	config := sql.Config{
		Driver: "postgres",
		ConnectionString: sql.ConnectionString{
			Host:     "localhost",
			Protocol: "postgres",
			Port:     "5432",
			Database: "connect_test",
			SSL:      sql.SSLInsecure,
			Credentials: sql.UsernamePasswordCredentials{
				Username: "dev_test",
				Password: "dev_test",
			},
		},
	}

	var err error
	s.conn, err = sql.NewConnection(s.ctx, config)
	s.Require().NoError(err)
	s.Require().NoError(s.conn.Ping())

	// Create test table
	_, err = s.conn.Exec(s.ctx, `
		CREATE TABLE IF NOT EXISTS test_records (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		)
	`)
	s.Require().NoError(err)
}

func (s *ConnectionTestSuite) TearDownSuite() {
	_, err := s.conn.Exec(s.ctx, `DROP TABLE IF EXISTS test_records`)
	s.Require().NoError(err)
	s.Require().NoError(s.conn.Close())
}

func (s *ConnectionTestSuite) SetupTest() {
	// Clean up any existing data before each test
	_, err := s.conn.Exec(s.ctx, `DELETE FROM test_records`)
	s.Require().NoError(err)
}

func (s *ConnectionTestSuite) TestNamedGetWithDeleteReturning() {
	// Insert test data
	_, err := s.conn.Exec(s.ctx, `INSERT INTO test_records (name) VALUES ($1)`, "test1")
	s.Require().NoError(err)

	// Test 1: Delete existing record with RETURNING
	var result testRecord
	err = s.conn.NamedGet(s.ctx, &result, `
		DELETE FROM test_records 
		WHERE name = :name 
		RETURNING id, name
	`, map[string]interface{}{
		"name": "test1",
	})
	s.Require().NoError(err)
	s.Equal("test1", result.Name)

	// Test 2: Delete non-existent record with RETURNING
	err = s.conn.NamedGet(s.ctx, &result, `
		DELETE FROM test_records 
		WHERE name = :name 
		RETURNING id, name
	`, map[string]interface{}{
		"name": "nonexistent",
	})
	s.ErrorIs(err, dbsql.ErrNoRows)
}

func (s *ConnectionTestSuite) TestNamedGetWithSelect() {
	// Insert test data
	_, err := s.conn.Exec(s.ctx, `INSERT INTO test_records (name) VALUES ($1)`, "test1")
	s.Require().NoError(err)

	// Test 1: Select existing record
	var result testRecord
	err = s.conn.NamedGet(s.ctx, &result, `
		SELECT id, name 
		FROM test_records 
		WHERE name = :name
	`, map[string]interface{}{
		"name": "test1",
	})
	s.Require().NoError(err)
	s.Equal("test1", result.Name)

	// Test 2: Select non-existent record
	err = s.conn.NamedGet(s.ctx, &result, `
		SELECT id, name 
		FROM test_records 
		WHERE name = :name
	`, map[string]interface{}{
		"name": "nonexistent",
	})
	s.ErrorIs(err, dbsql.ErrNoRows)
}

func (s *ConnectionTestSuite) TestNamedGetWithUpdateReturning() {
	// Insert test data
	_, err := s.conn.Exec(s.ctx, `INSERT INTO test_records (name) VALUES ($1)`, "test1")
	s.Require().NoError(err)

	// Test 1: Update existing record with RETURNING
	var result testRecord
	err = s.conn.NamedGet(s.ctx, &result, `
		UPDATE test_records 
		SET name = :new_name 
		WHERE name = :old_name 
		RETURNING id, name
	`, map[string]interface{}{
		"old_name": "test1",
		"new_name": "test2",
	})
	s.Require().NoError(err)
	s.Equal("test2", result.Name)

	// Test 2: Update non-existent record with RETURNING
	err = s.conn.NamedGet(s.ctx, &result, `
		UPDATE test_records 
		SET name = :new_name 
		WHERE name = :old_name 
		RETURNING id, name
	`, map[string]interface{}{
		"old_name": "nonexistent",
		"new_name": "test3",
	})
	s.ErrorIs(err, dbsql.ErrNoRows)
}

func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionTestSuite))
}
