package storage

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PgConn struct {
	Host    string
	Port    string
	User    string
	Passwd  string
	Schema  string
	Options *string
	Logger  *zap.Logger
}

func NewPgPool(pg PgConn) *pgxpool.Pool {
	url := "postgres://{{username}}:{{password}}@{{host}}:{{port}}/{{schema}}"
	url = strings.Replace(url, "{{host}}", pg.Host, -1)
	url = strings.Replace(url, "{{port}}", pg.Port, -1)
	url = strings.Replace(url, "{{username}}", pg.User, -1)
	url = strings.Replace(url, "{{password}}", pg.Passwd, -1)
	url = strings.Replace(url, "{{schema}}", pg.Schema, -1)
	if pg.Options != nil {
		url += "?" + *pg.Options
	}
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		pg.Logger.Error("failed to settle postgres connection", zap.Error(err))
		panic(err)
	}
	return pool
}
