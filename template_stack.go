package solislog

type stack[T any] struct {
	items []T
}

func newStack[T any](capacity int) *stack[T] {
	return &stack[T]{
		items: make([]T, 0, capacity),
	}
}

func (s *stack[T]) push(item T) {
	s.items = append(s.items, item)
}

func (s *stack[T]) pop() (T, bool) {
	var zero T
	if s.isEmpty() {
		return zero, false
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items[index] = zero
	s.items = s.items[:index]

	return item, true
}

func (s *stack[T]) peek() (T, bool) {
	if s.isEmpty() {
		var zero T
		return zero, false
	}

	return s.items[len(s.items)-1], true
}

func (s *stack[T]) len() int {
	return len(s.items)
}

func (s *stack[T]) isEmpty() bool {
	return len(s.items) == 0
}

func (s *stack[T]) clear() {
	for i := range s.items {
		var zero T
		s.items[i] = zero
	}
	s.items = s.items[:0]
}
