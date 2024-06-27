package structure

type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Len() int {
	return len(q.items)
}

func (q *Queue[T]) Add(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Drop() T {
	if q.Len() == 0 {
		var zero T
		return zero
	}
	item := q.items[0]
	q.items = q.items[1:]

	return item
}
