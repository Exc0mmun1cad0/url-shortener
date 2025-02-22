package postgres

import (
	"fmt"
	"net/url"
)

type connStr string

func FormConnStr(cfg Config) connStr {

	vals := url.Values{}
	vals.Set("sslmode", cfg.SSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:     cfg.DBName,
		RawQuery: vals.Encode(),
	}

	return connStr(u.String())
}

func (c connStr) WithMigrationsTable(migrationsTable string) connStr {
	u, err := url.Parse(string(c))
	if err != nil {
		panic("cannot parse connection string")
	}

	q := u.Query()
	q.Set("x-migrations-table", migrationsTable)

	u.RawQuery = q.Encode()

	return connStr(u.String())
}

// For better readability of chain function calls
func (c connStr) String() string {
	return string(c)
}
