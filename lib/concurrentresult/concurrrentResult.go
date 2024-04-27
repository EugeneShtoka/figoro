package concurrentresult

import (
	"context"
)

type ConcurrentResult[T any] struct {
    resultsChannel chan T
	errorChanel chan error
	cancel context.CancelFunc
}

func New[T any](ctx context.Context) *ConcurrentResult[T] {
	resultsChannel := make(chan T)
	errorChanel := make(chan error)
	_, cancel := context.WithCancel(ctx)

	return &ConcurrentResult[T]{
		resultsChannel,
		errorChanel,
		cancel,
	}
}

func (this *ConcurrentResult[T]) SendResult(result T) {
	this.resultsChannel <- result
}

func (this *ConcurrentResult[T]) SendError(err error) {
	this.errorChanel <- err
}

func (this *ConcurrentResult[T]) Cancel() {
	this.cancel()
}

func (this *ConcurrentResult[T]) Results(resCount int) ([]T, error) {
	results := make([]T, resCount)
	for i := 0; i < resCount; i++ {
		select {
			case result := <-this.resultsChannel:
				results[i] = result
			case err := <-this.errorChanel:
				return make([]T, resCount), err
		}
	}
	return results, nil
}