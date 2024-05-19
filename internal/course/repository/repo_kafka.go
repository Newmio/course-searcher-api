package repository

import (
	"github.com/IBM/sarama"
	"github.com/Newmio/newm_helper"
)

type kafkaCourseRepo struct {
	kafka sarama.SyncProducer
}

func NewKafkaCourseRepo(kafka sarama.Client) IKafkaCourseRepo {
	p, err := sarama.NewSyncProducerFromClient(kafka)
	if err != nil {
		panic(err)
	}
	return &kafkaCourseRepo{kafka: p}
}

func (r *kafkaCourseRepo) CreateCourseEvent(value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: "create_course",
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := r.kafka.SendMessage(msg)
	if err != nil {
		return newm_helper.Trace(err)
	}

	return nil
}
