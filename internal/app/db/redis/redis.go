package redis

import (
	"context"
	"fmt"
	"searcher/internal/app/db"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func OpenDb() (*redis.Client, error) {
	v := viper.New()
	v.AddConfigPath("internal/app/db/redis")
	v.SetConfigName("config")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	db, err := initDb(db.Config{
		Host:     v.GetString("host"),
		Port:     v.GetString("port"),
		User:     v.GetString("user"),
		Password: v.GetString("password"),
		DbName:   v.GetString("dbName"),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDb(c db.Config) (*redis.Client, error) {
	dbName, err := strconv.Atoi(c.DbName)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Password: c.Password,
		DB:       dbName,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
