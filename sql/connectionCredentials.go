package sql

import (
	"context"

	"goeasy.dev/errors"
)

// UsernamePasswordCredentials authenticates with the database using a username and password
type UsernamePasswordCredentials struct {
	Username string
	Password string
}

func (u UsernamePasswordCredentials) GetCredentials(_ context.Context) (string, string, error) {
	return u.Username, u.Password, nil
}

type CredentialProvider interface {
	Username() (string, error)
	Password() (string, error)
}

type providerCredentials struct {
	provider CredentialProvider
}

func ProviderCredentials(provider CredentialProvider) Credentials {
	return providerCredentials{provider: provider}
}

func (p providerCredentials) GetCredentials(_ context.Context) (string, string, error) {
	username, err := p.provider.Username()
	if err != nil {
		return "", "", errors.Wrap(err, "unable to retrieve username from provider")
	}

	password, err := p.provider.Password()
	if err != nil {
		return "", "", errors.Wrap(err, "unable to retrieve password from provider")
	}

	return username, password, nil
}
