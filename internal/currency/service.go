package currency

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"phoenix/internal/broker"
	"phoenix/internal/common"
)

// CurrencyService — сервис валют, обменивается сообщениями с accounts через AMQP.
type CurrencyService struct {
	*common.BaseService
	broker      *broker.Broker
	instanceUID string
}

// NewCurrencyService создаёт и запускает экземпляр CurrencyService.
func NewCurrencyService(cfg *CurrencyConfig) (*CurrencyService, error) {
	cs := &CurrencyService{BaseService: common.NewBaseService()}
	cs.instanceUID = fmt.Sprintf("%p", cs)

	b, err := broker.New(cfg.Amqp.URL, cfg.Amqp.Exchange)
	if err != nil {
		return nil, err
	}
	cs.broker = b

	queueName := fmt.Sprintf("currency.%s", sanitizeQueuePart(cs.instanceUID))
	cs.StartBackgroundTask(func(ctx context.Context) {
		err := b.Subscribe(ctx, queueName, []string{broker.BindingCurrencyPing}, cs.handleAMQPMessage)
		if err != nil && ctx.Err() == nil {
			log.Printf("AMQP subscribe failed: %v", err)
		}
	})

	log.Printf("Currency service instance uid=%s queue=%s", cs.instanceUID, queueName)
	return cs, nil
}

func (cs *CurrencyService) handleAMQPMessage(ctx context.Context, routingKey string, body []byte) error {
	if ctx.Err() != nil {
		return nil
	}

	msgType, err := broker.MessageType(body)
	if err != nil {
		log.Printf("AMQP message ignored (routing_key=%s): %v", routingKey, err)
		return nil
	}

	if msgType != broker.TypePingCurrenciesSupportStatus {
		return nil
	}

	pongBody, err := broker.NewPongBody(cs.instanceUID, time.Now())
	if err != nil {
		log.Printf("Failed to build pong message: %v", err)
		return nil
	}

	if err := cs.broker.Publish(ctx, broker.RoutingPongCurrenciesSupportStatus, pongBody); err != nil {
		log.Printf("Failed to publish pong:CurrenciesSupportStatus: %v", err)
		return nil
	}

	log.Printf("Sent pong:CurrenciesSupportStatus uid=%s routing_key=%s", cs.instanceUID, broker.RoutingPongCurrenciesSupportStatus)
	return nil
}

// Shutdown останавливает сервис и закрывает AMQP-соединение.
func (cs *CurrencyService) Shutdown(ctx context.Context) error {
	err := cs.BaseService.Shutdown(ctx)
	if cs.broker != nil {
		if closeErr := cs.broker.Close(); closeErr != nil {
			if err != nil {
				return fmt.Errorf("%w; amqp close: %v", err, closeErr)
			}
			return closeErr
		}
	}
	return err
}

func sanitizeQueuePart(s string) string {
	replacer := strings.NewReplacer("0x", "", "*", "", " ", "", ".", "")
	return replacer.Replace(s)
}
