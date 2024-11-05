package kafka

import (
	"context"
	"fmt"
	"mpc/internal/infrastructure/config"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type Writer = kafka.Writer
type Reader = kafka.Reader

type kafkaOptions struct {
	Topic string
}

type KafkaOption func(*kafkaOptions)

func defaultKafkaOptions(cfg *config.KafkaConfig) kafkaOptions {
	return kafkaOptions{
		Topic: cfg.Topic,
	}
}

func WithTopic(topic string) KafkaOption {
	return func(o *kafkaOptions) {
		o.Topic = topic
	}
}

func NewKafkaProducer(cfg *config.KafkaConfig, opts ...KafkaOption) (*Writer, error) {
	options := defaultKafkaOptions(cfg)
	for _, opt := range opts {
		opt(&options)
	}
	fmt.Print("Run into this")

	err := createTopicIfNotExists(cfg.Brokers, options.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}
	fmt.Print("Run into this 2")
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Brokers,
		Topic:    options.Topic,
		Balancer: &kafka.LeastBytes{},
	})
	return producer, nil
}

func NewKafkaConsumer(cfg *config.KafkaConfig, opts ...KafkaOption) (*kafka.Reader, error) {
	options := defaultKafkaOptions(cfg)
	for _, opt := range opts {
		opt(&options)
	}

	if options.Topic != cfg.Topic {
		err := createTopicIfNotExists(cfg.Brokers, options.Topic)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic: %w", err)
		}
	}

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   options.Topic,
	})
	return consumer, nil
}

func CloseKafkaConnections(producer *kafka.Writer, consumers ...*kafka.Reader) error {
	var err error
	if producer != nil {
		if e := producer.Close(); e != nil {
			err = e
		}
	}
	for _, consumer := range consumers {
		if consumer != nil {
			if e := consumer.Close(); e != nil && err == nil {
				err = e
			}
		}
	}
	return err
}

func createTopicIfNotExists(brokers []string, topic string) error {
	// Add timeout for initial connection
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Connect to the first broker
	conn, err := dialer.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker %s: %w", brokers[0], err)
	}
	defer conn.Close()

	// Add timeout for controller request
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := dialer.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka controller at %s: %w", controllerAddr, err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	// Add timeout for create topics request
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		if err == kafka.TopicAlreadyExists {
			return nil // Topic exists, not an error
		}
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}

	return nil
}

func PublishMessage(producer *kafka.Writer, key, value []byte) error {
	return producer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)
}

func ConsumeMessages(consumer *kafka.Reader, handler func(kafka.Message) error) error {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		if err := handler(msg); err != nil {
			return err
		}
	}
}

// ReadNewMessage reads only new messages published after this function is called
func ReadNewMessage(ctx context.Context, reader *kafka.Reader) (kafka.Message, error) {
	// First, seek to the end of the topic
	err := reader.SetOffset(kafka.LastOffset)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("failed to set offset: %w", err)
	}

	// Now read the next message, which should be a new one
	return reader.ReadMessage(ctx)
}
