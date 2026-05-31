package apiserver

import (
	"context"
	"errors"
	"net/http"
	"phoenix/internal/config"
	"phoenix/pkg/types"
	"time"
)

// ApiServer HTTP-сервер предоставляет API.
type ApiServer struct {
	addr types.Address
	srv  *http.Server
}

// CreateAndStart создает и запускает в гороутине экземпляр ApiServer.
func CreateAndStart(cfg config.AccountsConfig, handler http.Handler) (*ApiServer, error) {
	address := types.Address{Host: cfg.Address.Host, Port: cfg.Address.Port}
	apiServer := &ApiServer{
		addr: address,
		srv: &http.Server{
			Addr:         address.String(),
			Handler:      handler,
			ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,  // Таймаут на чтение запроса
			WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second, // Таймаут на запись ответа
		},
	}

	errCh := make(chan error, 1)
	// Запуск сервера в отдельной гороутине
	go func() {
		if err := apiServer.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// Ожидаю 500 ms проверяя что сервер не вернул ошибку при старте
	select {
	case err := <-errCh:
		return nil, err
	case <-time.After(500 * time.Millisecond):
		return apiServer, nil
	}
}

// Address адрес на котором сервер ожидает входящие подключения.
func (as ApiServer) Address() types.Address {
	return as.addr
}

// Shutdown останавливает и освобождает ресурсы запущенного экземпляра ApiServer.
func (as *ApiServer) Shutdown(ctx context.Context) error {
	return as.srv.Shutdown(ctx)
}
