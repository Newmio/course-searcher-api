package kafka

import (
	"fmt"
	"searcher/internal/app/storage"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func OpenBroker() (sarama.Client, error) {
	v := viper.New()
	v.AddConfigPath("internal/app/storage/kafka")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return initBroker(storage.Config{
		Host: v.GetString("host"),
		Port: v.GetString("port"),
	})
}

func initBroker(c storage.Config) (sarama.Client, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	return sarama.NewClient([]string{fmt.Sprintf("%s:%s", c.Host, c.Port)}, config)
}
