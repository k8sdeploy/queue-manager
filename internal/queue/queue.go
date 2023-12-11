package queue

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	ConfigBuilder "github.com/keloran/go-config"
	"github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	ConfigBuilder.Config
}

func NewQueue(cfg ConfigBuilder.Config) *Queue {
	return &Queue{
		Config: cfg,
	}
}

func (q *Queue) Start() error {
	conString := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", q.Rabbit.Username, q.Rabbit.Password, q.Rabbit.Host, q.Rabbit.Port, q.Rabbit.VHost)
	if q.Rabbit.Port == 0 {
		conString = fmt.Sprintf("amqp://%s:%s@%s/%s", q.Rabbit.Username, q.Rabbit.Password, q.Rabbit.Host, q.Rabbit.VHost)
	}

	conn, err := amqp091.Dial(conString)
	if err != nil {
		return logs.Errorf("queue: unable to dial: %v, %s", err, conString)
	}

	fmt.Sprintf("%+v", conn)

	return nil
}
