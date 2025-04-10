package sql

import (
	"context"
	"fmt"

	"goeasy.dev/errors"
	"goeasy.dev/observability/log"
)

// ConnectionString is used to build the DSN for connecting to the database
// Raw can be set to specify the DSN directly
type ConnectionString struct {
	Raw         string
	Protocol    string
	Host        string
	Port        string
	SSL         ConnectionSSL
	Credentials Credentials
	Database    string
}

// Credentials are used to retrieve the credentials for connection
type Credentials interface {
	GetCredentials(context.Context) (username string, password string, err error)
}

// ConnectionSSL is used to specify the SSL options
type ConnectionSSL interface {
	GetSSLOptions() map[string]string
}

type simpleSSL map[string]string

func (s simpleSSL) GetSSLOptions() map[string]string {
	return s
}

var SSLInsecure simpleSSL = simpleSSL{
	"sslmode": "disable",
}

// Build the DSN for connection
func (c *ConnectionString) Build(ctx context.Context) (string, error) {
	log.Debug("sql: building connection string")
	if c.Raw != "" {
		return c.Raw, nil
	}

	username, password, err := c.Credentials.GetCredentials(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve database credentials")
	}

	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s/%s", c.Protocol, username, password, c.Host, c.Port, c.Database)
	options := c.SSL.GetSSLOptions()
	if len(options) != 0 {
		connectionString += "?" + buildOptions(options)
	}

	return connectionString, nil
}

func buildOptions(options map[string]string) string {
	out := ""
	i := 0
	j := len(options)
	for k, v := range options {
		out += fmt.Sprintf("%s=%s", k, v)
		if i != j-1 {
			out += "&"
		}

		i++
	}

	return out
}
