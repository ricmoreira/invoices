package services

import (
	"encoding/json"
	"invoices/config"
	"invoices/models/request"
	"invoices/models/response"
	"invoices/util/errors"
	"log"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer allows sending messages to Kafka broker
type KafkaProducer struct {
	prod   *kafka.Producer
	config *config.Config
}

// NewKafkaProducer is the KafkaProducer constructor
func NewKafkaProducer(config *config.Config) (*KafkaProducer, error) {
	kp := KafkaProducer{
		config: config,
	}

	err := kp.Connect()
	if err != nil {
		return nil, err
	}

	return &kp, nil
}

// Connect connects instance to Kafka server enabling instance to send messages
func (kp *KafkaProducer) Connect() error {
	config := &kafka.ConfigMap{
		"bootstrap.servers":    kp.config.KafkaConsumerConfig.BootstrapServers,
		"client.id":            "invoices",
		"default.topic.config": kafka.ConfigMap{"acks": "all"},
		"message.max.bytes":    kp.config.MessageMaxBytes,
	}

	p, err := kafka.NewProducer(config)

	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}

	kp.prod = p
	return nil
}

// SendInvoiceToTopic sends an invoice to provided topic name
func (kp *KafkaProducer) SendInvoiceToTopic(topic string, request *mrequest.InvoiceCreate) *mresponse.ErrorResponse {

	jptBytes, err := json.Marshal(request)
	if err != nil {
		e := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
		return e
	}

	// don't allow to process more than max size allowed of data
	if len(jptBytes) > kp.config.MessageMaxBytes {
		s := strconv.Itoa(kp.config.MessageMaxBytes)
		e := errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, "The invoice has to many data to processs. Max allowed: "+s+" bytes of data.")
		return e
	}

	e := kp.sendMessage(topic, jptBytes)

	if e != nil {
		return e
	}

	return nil
}

// SendMessage sends a message to Kafka to provided topic
// errors.HandleErrorResponse(500, errors.SERVICE_UNAVAILABLE, errors.DefaultErrorsMessages[errors.SERVICE_UNAVAILABLE], nil)
func (kp *KafkaProducer) sendMessage(topic string, message []byte) *mresponse.ErrorResponse {

	deliveryChan := make(chan kafka.Event, 10000)

	mes := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}

	err := kp.prod.Produce(&mes, deliveryChan)

	if err != nil {
		return errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, err.Error())
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		close(deliveryChan)
		return errors.HandleErrorResponse(errors.SERVICE_UNAVAILABLE, nil, m.TopicPartition.Error.Error())
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		close(deliveryChan)
		return nil
	}
}
