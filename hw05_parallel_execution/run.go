package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		n = 1
	}

	// Канал для передачи задач воркерам.
	tasksCh := make(chan Task)

	// определение группы горутин и создание счётчика.
	var wg sync.WaitGroup
	var errorsCount int64

	maxErrors := int64(m)  // максимально допустимое колличество ошибок
	ignoreErrors := m <= 0 // игнорирование ошибок

	// Запуск n воркеров.
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range tasksCh {
				if err := task(); err != nil {
					if !ignoreErrors && atomic.AddInt64(&errorsCount, 1) >= maxErrors {
						return
					}
				}
			}
		}()
	}

	// Отправка задач в канал.
	// Если лимит ошибок достигнут — прекращаем отправку новых задач.
	for _, task := range tasks {
		if !ignoreErrors && atomic.LoadInt64(&errorsCount) >= maxErrors {
			break
		}
		tasksCh <- task
	}

	// Закрываем канал, сигнализируя воркерам об окончании задач. Ждём завершения всех воркеров.
	close(tasksCh)
	wg.Wait()

	// Если лимит ошибок был превышен — возвращаем ошибку.
	if !ignoreErrors && atomic.LoadInt64(&errorsCount) >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
