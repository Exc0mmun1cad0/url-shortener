package postgres

import (
	"fmt"
	"net/url"
	"url-shortener/internal/config"
)

// From connection string forms dataSource for postgresql using url package.
// It can form connection string from either Config struct or separate parameters.
// Params: user, password, host, port, dbName, sslMode, migrationsTable
func FormConnStr(cfg config.Postgres, params ...string) string {
	q := url.Values{}
	var dataSrc url.URL

	if cfg != (config.Postgres{}) {
		q.Set("sslmode", cfg.SSLMode)
		dataSrc = url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(cfg.User, cfg.Password),
			Host:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Path:     cfg.DBName,
			RawQuery: q.Encode(),
		}
	} else {
		if len(params) < 7 {
			panic("not enough arguments specified for postgres connection")
		}

		q.Set("x-migrations-table", params[6])
		q.Set("sslmode", params[5])
		dataSrc = url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(params[2], params[3]),
			Host:     fmt.Sprintf("%s:%s", params[0], params[1]),
			Path:     params[4],
			RawQuery: q.Encode(),
		}
	}

	return dataSrc.String()
}
