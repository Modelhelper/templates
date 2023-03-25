package mssql

import (
	"context"
	"database/sql"
	"{{ .Name | kebab }}/service"
)

type SqlConnectionManager interface {
	Close() error
	Get(ctx context.Context) (*sql.Conn, error)
}

func NewSqlConnectionManager(config *service.SqlConfig) SqlConnectionManager {
	return &sqlConnManager{config: config}
}

type sqlConnManager struct {
	pool   *sql.DB
	config *service.SqlConfig
}

func (s *sqlConnManager) Close() error {
	if s.pool != nil {
		return s.pool.Close()
	}
	return nil
}

func (s *sqlConnManager) Get(ctx context.Context) (*sql.Conn, error) {
	var err error

	if s.pool == nil {
		s.pool, err = sql.Open("sqlserver", s.config.ConnectionString)
		if err != nil {
			return nil, err
		}
	}

	return s.pool.Conn(ctx)
}
