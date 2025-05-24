package redis

import (
	"context"
	"fmt"
	"searcher/internal/app/storage"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func OpenDb() (*redis.Client, error) {
	// v := viper.New()
	// v.AddConfigPath("internal/app/storage/redis")
	// v.SetConfigName("config")

	// err := v.ReadInConfig()
	// if err != nil {
	// 	return nil, err
	// }

	return initDb(storage.Config{
		Host:     "localhost",
		Port:     "6379",
		User:     "",
		Password: "qwerty",
		DbName:   "0",
	})

	// return initDb(storage.Config{
	// 	Host:     v.GetString("host"),
	// 	Port:     v.GetString("port"),
	// 	User:     v.GetString("user"),
	// 	Password: v.GetString("password"),
	// 	DbName:   v.GetString("dbName"),
	// })
}

func initDb(c storage.Config) (*redis.Client, error) {
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
