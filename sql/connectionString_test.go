package sql_test

import (
	"context"
	"testing"

	"goeasy.dev/sql"

	"github.com/stretchr/testify/assert"
)

type testProvider struct{}

func (testProvider) Username() (string, error) {
	return "testProvider", nil
}

func (testProvider) Password() (string, error) {
	return "testPassword", nil
}

func TestConnectionString(t *testing.T) {
	testCases := []struct {
		desc             string
		expected         string
		connectionString sql.ConnectionString
		errFunc          func(assert.TestingT, error, ...interface{}) bool
	}{
		{
			desc: "UserPass",
			connectionString: sql.ConnectionString{
				Host:     "localhost",
				Protocol: "postgres",
				Port:     "5432",
				Database: "testdb",
				SSL:      sql.SSLInsecure,
				Credentials: sql.UsernamePasswordCredentials{
					Username: "test",
					Password: "testing",
				},
			},
			expected: "postgres://test:testing@localhost:5432/testdb?sslmode=disable",
			errFunc:  assert.NoError,
		},
		{
			desc: "Provider",
			connectionString: sql.ConnectionString{
				Host:        "localhost",
				Protocol:    "postgres",
				Port:        "5432",
				Database:    "testdb",
				SSL:         sql.SSLInsecure,
				Credentials: sql.ProviderCredentials(testProvider{}),
			},
			expected: "postgres://testProvider:testPassword@localhost:5432/testdb?sslmode=disable",
			errFunc:  assert.NoError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			connString, err := tC.connectionString.Build(context.Background())
			if !tC.errFunc(t, err) {
				return
			}

			assert.Equal(t, tC.expected, connString)
		})
	}
}
