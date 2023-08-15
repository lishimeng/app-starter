package app

import (
	"github.com/lishimeng/app-starter/amqp"
	"github.com/lishimeng/app-starter/amqp/rabbit"
)

func AmqpPub(m rabbit.Message) error {
	return amqp.Publish(GetAmqp(), m)
}

func MqttPub(topic string, qos byte, retained bool, data []byte) error {
	return GetMqtt().Publish(topic, qos, retained, data)
}

func MqttPublish(topic string, data []byte) error {
	return GetMqtt().Publish(topic, 0, false, data)
}
