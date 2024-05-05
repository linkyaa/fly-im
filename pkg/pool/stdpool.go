package pool

import "sync"

type (
	stdPool[T any] struct {
		pool sync.Pool
	}
)

func (s *stdPool[T]) Get() T {
	return s.pool.Get().(T)
}

func (s *stdPool[T]) Put(t T) {
	s.pool.Put(t)
}

func NewStdPool[T any](factory func() T) Pooler[T] {
	res := &stdPool[T]{
		pool: sync.Pool{New: func() interface{} { return factory() }},
	}
	return res
}
