package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Rdb struct {
	User     string
	Password string
	Host     string
	Port     string
}

func NewRdb(user, password, host, port string) *Rdb {
	return &Rdb{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (r *Rdb) Connect() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Username: r.User,
		Password: r.Password,
	})

	return redis
}
