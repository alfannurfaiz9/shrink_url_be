package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PsqlDB struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewPsqlDB(user, password, host, port, dbName string) *PsqlDB {
	return &PsqlDB{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		DBName:   dbName,
	}
}

func (db *PsqlDB) Connect() (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.DBName)

	return pgxpool.New(context.Background(), connString)
}
