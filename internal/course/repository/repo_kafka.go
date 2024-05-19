package repository

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/Newmio/newm_helper"
)

type kafkaCourseRepo struct {
	kafka sarama.SyncProducer
	admin sarama.ClusterAdmin
}

func NewKafkaCourseRepo(kafka sarama.Client) IKafkaCourseRepo {
	p, err := sarama.NewSyncProducerFromClient(kafka)
	if err != nil {
		panic(err)
	}

	admin, err := sarama.NewClusterAdminFromClient(kafka)
	if err != nil {
		panic(err)
	}
	return &kafkaCourseRepo{kafka: p, admin: admin}
}

func (r *kafkaCourseRepo) CreateCourseEvent(value []byte) error {
	topic := fmt.Sprintf("create_course_%d", time.Now().Day())

	if err := r.createTopicIfNotExists(topic); err != nil {
		return newm_helper.Trace(err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := r.kafka.SendMessage(msg)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}

func (r *kafkaCourseRepo) createTopicIfNotExists(topic string) error {
	var exists bool
	dateDell := "86400000"

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1, // Количество партиций
		ReplicationFactor: 1, // Фактор репликации
		ConfigEntries: map[string]*string{
			"retention.ms": &dateDell, // Время жизни топика в миллисекундах (здесь 24 часа)
		},
	}

	topictList, err := r.admin.ListTopics()
	if err != nil {
		return newm_helper.Trace(err)
	}

	for name := range topictList {
		if name == topic {
			exists = true
		}
	}

	if !exists {
		err = r.admin.CreateTopic(topic, topicDetail, false)
		if err != nil {
			return newm_helper.Trace(err)
		}
	}

	return nil
}
