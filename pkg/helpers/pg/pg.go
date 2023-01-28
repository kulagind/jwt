package pg

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	Timeout  int
}

func NewPoolConfig(cfg *Config) (*pgxpool.Config, error) {
	connectionConfig := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DbName,
		cfg.Timeout,
	)

	poolConfig, err := pgxpool.ParseConfig(connectionConfig)
	if err != nil {
		return nil, err
	}

	return poolConfig, nil
}

func NewConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CheckSqlError(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) || err.Error() == code {
		return true
	}
	return false
}
