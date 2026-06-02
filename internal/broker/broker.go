package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler обрабатывает входящее AMQP-сообщение.
type MessageHandler func(ctx context.Context, routingKey string, body []byte) error

// Broker — клиент topic exchange RabbitMQ.
type Broker struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

// New подключается к RabbitMQ и объявляет topic exchange.
func New(amqpURL, exchange string) (*Broker, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("amqp dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("declare exchange %q: %w", exchange, err)
	}

	return &Broker{conn: conn, ch: ch, exchange: exchange}, nil
}

// Publish публикует сообщение в topic exchange.
func (b *Broker) Publish(ctx context.Context, routingKey string, body []byte) error {
	return b.ch.PublishWithContext(ctx, b.exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// Subscribe объявляет очередь, привязывает её к exchange и читает сообщения до отмены ctx.
func (b *Broker) Subscribe(ctx context.Context, queueName string, bindingKeys []string, handler MessageHandler) error {
	q, err := b.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue %q: %w", queueName, err)
	}

	for _, key := range bindingKeys {
		if err := b.ch.QueueBind(q.Name, key, b.exchange, false, nil); err != nil {
			return fmt.Errorf("bind queue %q to %q: %w", q.Name, key, err)
		}
	}

	deliveries, err := b.ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume queue %q: %w", q.Name, err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-deliveries:
			if !ok {
				return nil
			}
			if err := handler(ctx, d.RoutingKey, d.Body); err != nil {
				return err
			}
		}
	}
}

// Close закрывает канал и соединение с брокером.
func (b *Broker) Close() error {
	if b.ch != nil {
		_ = b.ch.Close()
	}
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}
