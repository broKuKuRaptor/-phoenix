package accounts

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// updateCurrencySupportState задача обновления CurrencySupportState
func updateCurrencySupportState(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	defer wg.Done() // Уменьшение счетчика работающих горутин

	ticker := time.NewTicker(interval)
	defer ticker.Stop() // Освобождение ресурсов тикера

	for {
		select {
		case <-ctx.Done():
			return // Получен сигнал остановки (контекст отменен)
		case <-ticker.C:
			// Наступил интервал времени
			fmt.Println("Hello, World")
			// Место для запуска переодически выполняемой задачи
		}
	}

}
