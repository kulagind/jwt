package pg

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
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

type CustomSqlConn struct {
	Exec     func(ctx context.Context, sql string, args ...any) (interface{}, error)
	QueryRow func(ctx context.Context, sql string, args ...any) *CustomRow
	Query    func(ctx context.Context, sql string, args ...any) (*CustomRows, error)
	Mock     *pgxpoolmock.MockPgxPool
	Close    func()
}

type CustomRow struct {
	Scan func(dest ...any) error
}

type CustomRows struct {
	Close func()
	Next  func() bool
	Err   func() error
	Scan  func(dest ...any) error
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

func NewConnection(poolConfig *pgxpool.Config) (*CustomSqlConn, error) {
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return &CustomSqlConn{
		Exec: func(ctx context.Context, sql string, args ...any) (interface{}, error) {
			return conn.Exec(ctx, sql, args...)
		},
		Query: func(ctx context.Context, sql string, args ...any) (*CustomRows, error) {
			pgxrows, err := conn.Query(ctx, sql, args...)
			if err != nil {
				return nil, err
			}
			return &CustomRows{
				Close: pgxrows.Close,
				Next:  pgxrows.Next,
				Err:   pgxrows.Err,
				Scan:  pgxrows.Scan,
			}, nil
		},
		QueryRow: func(ctx context.Context, sql string, args ...any) *CustomRow {
			return &CustomRow{
				Scan: conn.QueryRow(ctx, sql, args...).Scan,
			}
		},
	}, nil
}

func CheckSqlError(err error, code string) bool {
	var pgErr *pgconn.PgError
	if err == nil {
		return false
	}
	if errors.As(err, &pgErr) || err.Error() == code {
		return true
	}
	return false
}

func NewMockConnection(t *testing.T) (*CustomSqlConn, error) {
	ctrl := gomock.NewController(t)

	mock := pgxpoolmock.NewMockPgxPool(ctrl)

	return &CustomSqlConn{
		Exec: func(ctx context.Context, sql string, args ...any) (interface{}, error) {
			return mock.Exec(ctx, sql, args...)
		},
		Query: func(ctx context.Context, sql string, args ...any) (*CustomRows, error) {
			rows, err := mock.Query(ctx, sql, args...)
			if err != nil {
				return nil, err
			}
			return &CustomRows{
				Close: rows.Close,
				Next:  rows.Next,
				Err:   rows.Err,
				Scan:  rows.Scan,
			}, nil
		},
		QueryRow: func(ctx context.Context, sql string, args ...any) *CustomRow {
			return &CustomRow{
				Scan: mock.QueryRow(ctx, sql, args...).Scan,
			}
		},
		Mock: mock,
		Close: func() {
			ctrl.Finish()
		},
	}, nil
}
