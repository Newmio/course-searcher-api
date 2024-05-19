package postgres

import (
	"fmt"
	"searcher/internal/app/storage"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func OpenDb() (*sqlx.DB, error) {
	v := viper.New()
	v.AddConfigPath("internal/app/storage/postgres")
	v.SetConfigName("config")
	
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return initDb(storage.Config{
		Host:     v.GetString("host"),
		Port:     v.GetString("port"),
		User:     v.GetString("user"),
		Password: v.GetString("password"),
		DbName:   v.GetString("dbName"),
		SslMode:  v.GetString("sslMode"),
	})
}

func initDb(c storage.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host, c.Port, c.User, c.DbName, c.Password, c.SslMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
