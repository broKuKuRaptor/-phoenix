package common

import (
	"context"
	"log"
	"sync"
	"time"
)

// BaseService — обобщает базовые методы сервисов
type BaseService struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewBaseService — инициализирует базовый сервис и возвращает его контекст.
func NewBaseService() *BaseService {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseService{ctx: ctx, cancel: cancel}
}

// startPeriodicTask — запускает периодическую задачу
func (b *BaseService) StartPeriodicTask(periodic_task func(), interval time.Duration) {
	b.wg.Add(1)
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		defer b.wg.Done()

		// Защита приложения от паники в фоновом потоке
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Background task failed: %v\n", r)
			}
		}()

		// Цикл переодического вызова задачи
		for {
			select {
			case <-b.ctx.Done():
				return // Контекст отменен, выход
			case <-ticker.C:
				periodic_task() // Вызов периодической задачи
			}
		}
	}()
}

// Shutdown останавливает и освобождает ресурсы запущенного экземпляра BaseService.
func (b *BaseService) Shutdown(ctx context.Context) error {
	b.cancel() // Отмена контекста фоновых задач (сигнал к их завершению)
	// Канал закроется, когда завершаться все связанные с сервисом фоновые задачи.
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()
	// Завершение по окончанию фоновых задач, либо закрытию контекста
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
