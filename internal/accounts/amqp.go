package accounts

import (
	"context"
	"fmt"
	"log"

	"phoenix/internal/broker"
)

const currenciesSupportPingInterval = 3

func (as *AccountsService) queueName() string {
	return "accounts.pong"
}

func (as *AccountsService) startAMQP(cfg *AccountsConfig) error {
	b, err := broker.New(cfg.Amqp.URL, cfg.Amqp.Exchange)
	if err != nil {
		return err
	}
	as.broker = b

	as.StartBackgroundTask(func(ctx context.Context) {
		err := b.Subscribe(ctx, as.queueName(), []string{broker.BindingAccountsPong}, as.handleAMQPMessage)
		if err != nil && ctx.Err() == nil {
			log.Printf("AMQP subscribe failed: %v", err)
		}
	})

	return nil
}

func (as *AccountsService) handleAMQPMessage(ctx context.Context, routingKey string, body []byte) error {
	if ctx.Err() != nil {
		return nil
	}

	msgType, err := broker.MessageType(body)
	if err != nil {
		log.Printf("AMQP message ignored (routing_key=%s): %v", routingKey, err)
		return nil
	}

	if msgType != broker.TypePongCurrenciesSupportStatus {
		return nil
	}

	pong, err := broker.ParsePongBody(body)
	if err != nil {
		log.Printf("AMQP pong parse failed (routing_key=%s): %v", routingKey, err)
		return nil
	}

	log.Printf("Received pong:CurrenciesSupportStatus uid=%s ts=%s routing_key=%s", pong.UID, pong.TS, routingKey)
	return nil
}

func (as *AccountsService) updateCurrenciesSupportStatus() {
	if as.broker == nil {
		return
	}

	body, err := broker.NewPingBody()
	if err != nil {
		log.Printf("Failed to build ping message: %v", err)
		return
	}

	if err := as.broker.Publish(context.Background(), broker.RoutingPingCurrenciesSupportStatus, body); err != nil {
		log.Printf("Failed to publish ping:CurrenciesSupportStatus: %v", err)
	}
}

// Shutdown останавливает сервис и закрывает AMQP-соединение.
func (as *AccountsService) Shutdown(ctx context.Context) error {
	err := as.BaseService.Shutdown(ctx)
	if as.broker != nil {
		if closeErr := as.broker.Close(); closeErr != nil {
			if err != nil {
				return fmt.Errorf("%w; amqp close: %v", err, closeErr)
			}
			return closeErr
		}
	}
	return err
}
