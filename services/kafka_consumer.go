package services

import (
	"encoding/json"
	"invoices/config"
	"invoices/models/request"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	config      *config.Config
	invoiceServ *InvoiceService
	kafkaProducer *KafkaProducer
}

func NewKafkaConsumer(config *config.Config, ps *InvoiceService, kp *KafkaProducer) *KafkaConsumer {
	return &KafkaConsumer{
		config:      config,
		invoiceServ: ps,
		kafkaProducer: kp,
	}
}

func (kc *KafkaConsumer) Run() {

	log.Println("Start receiving from Kafka")

	configConsumer := kafka.ConfigMap{
		"bootstrap.servers":       kc.config.BootstrapServers,
		"group.id":                kc.config.GroupID,
		"auto.offset.reset":       kc.config.AutoOffsetReset,
		"auto.commit.enable":      kc.config.AutoCommitEnable,
		"auto.commit.interval.ms": kc.config.AutoCommitInterval,
	}

	c, err := kafka.NewConsumer(&configConsumer)

	if err != nil {
		panic(err)
	}

	topicsSubs := kc.config.TopicsSubscribed
	err = c.SubscribeTopics(topicsSubs, nil)

	if err != nil {
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)

		if err == nil {

			topic := *msg.TopicPartition.Topic

			switch topic {
			case "invoices":
				log.Println(`Reading an invoices topic message`)
				invoices, err := kc.parseInvoicesMessage(msg.Value)
				if err != nil {
					log.Printf("Error parsing event message value. Message %v \n Error: %s\n", msg.Value, err.Error())
					break
				}

				// save each invoice to database and produce message to kafka
				for _, invoice := range *invoices {
					inv, e := kc.invoiceServ.CreateOne(invoice) // save invoice to database

					log.Printf("%v",invoice)
					log.Printf("%v",inv)
					if e != nil {
						log.Printf("Error saving invoice to database\n Error: %s\n", e.Response)
					} else {
						kc.kafkaProducer.SendInvoiceToTopic("invoice_created", inv) // inform Kafka invoice created
					}
				}
			default: //ignore any other topics
			}
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func (kc *KafkaConsumer) parseInvoicesMessage(messageValue []byte) (*[]*mrequest.InvoiceCreate, error) {
	invoices := make([]*mrequest.InvoiceCreate, 0)
	err := json.Unmarshal(messageValue, &invoices)

	if err != nil {
		return nil, err
	}

	return &invoices, nil
}
