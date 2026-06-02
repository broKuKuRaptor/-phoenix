package broker

import (
	"encoding/json"
	"fmt"
	"time"
)

// PingMessage тело ping:CurrenciesSupportStatus.
type PingMessage struct {
	Type string `json:"type"`
}

// PongMessage тело pong:CurrenciesSupportStatus.
type PongMessage struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
	TS   string `json:"ts"`
}

// NewPingBody формирует JSON для ping:CurrenciesSupportStatus.
func NewPingBody() ([]byte, error) {
	return json.Marshal(PingMessage{Type: TypePingCurrenciesSupportStatus})
}

// NewPongBody формирует JSON для pong:CurrenciesSupportStatus.
func NewPongBody(uid string, ts time.Time) ([]byte, error) {
	return json.Marshal(PongMessage{
		Type: TypePongCurrenciesSupportStatus,
		UID:  uid,
		TS:   ts.UTC().Format(time.RFC3339),
	})
}

// ParsePongBody разбирает тело pong-сообщения.
func ParsePongBody(body []byte) (PongMessage, error) {
	var msg PongMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return PongMessage{}, fmt.Errorf("decode pong message: %w", err)
	}
	if msg.Type != TypePongCurrenciesSupportStatus {
		return PongMessage{}, fmt.Errorf("unexpected message type: %s", msg.Type)
	}
	return msg, nil
}

// MessageType возвращает поле type из произвольного JSON-сообщения.
func MessageType(body []byte) (string, error) {
	var envelope struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return "", fmt.Errorf("decode message envelope: %w", err)
	}
	return envelope.Type, nil
}
