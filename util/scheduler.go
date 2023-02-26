package util

import (
	"time"
)

type Promise[T any] struct {
	val   T
	err   error
	ready chan struct{}
}

type Func[T, V any] func(T) (time.Duration, V, error)

func (p *Promise[V]) Get() (V, error) {
	<-p.ready
	return p.val, p.err
}

func Schedule[T, V any](initDuration time.Duration, t T, f Func[T, V]) *Promise[V] {
	ready := make(chan struct{})
	p := Promise[V]{
		ready: ready,
	}

	ticker := time.NewTicker(initDuration)
	go func() {
		for {
			<-ticker.C
			var nextDuration time.Duration
			nextDuration, p.val, p.err = f(t)
			ticker.Reset(nextDuration)
			ready <- struct{}{}
		}
	}()

	return &p
}
