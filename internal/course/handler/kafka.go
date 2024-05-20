package handler

import (
	"fmt"
	"searcher/internal/course/service"
	"time"

	"github.com/IBM/sarama"
	"github.com/Newmio/newm_helper"
)

type KafkaHandler struct {
	s             service.ICourseService
	consumer      sarama.Consumer
	offsetManager sarama.OffsetManager
}

func NewKafkaHandler(client sarama.Client, s service.ICourseService) (*KafkaHandler, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	manager, err := sarama.NewOffsetManagerFromClient("course_manager", client)
	if err != nil {
		return nil, newm_helper.Trace(err)
	}

	return &KafkaHandler{
		s:             s,
		consumer:      consumer,
		offsetManager: manager,
	}, nil
}

func (h *KafkaHandler) Run(){
	go h.GetCourseEvent()
}

func (h *KafkaHandler) GetCourseEvent() {
	topic := fmt.Sprintf("course_event_%d", time.Now().Day())

	portition, err := h.consumer.ConsumePartition(topic, 0, 0)
	if err != nil {
		return
	}
	defer portition.Close()

	manager, err := h.offsetManager.ManagePartition(topic, 0)
	if err != nil {
		return
	}
	defer manager.Close()

	for {
		select {
		case msg := <-portition.Messages():

			flag, err := h.s.CheckExistsEventOffset(int(msg.Offset))
			if err != nil {
				return
			}

			if flag {
				BroadcastCourseEvent(msg.Value)
			}

		case <-time.After(time.Second * 60):
			manager.ResetOffset(0, "")
		}
	}
}
