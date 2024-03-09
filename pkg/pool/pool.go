package pool

/*
内存池
*/

type (
	Pooler[T any] interface {
		Get() T
		Put(T)
	}
)
