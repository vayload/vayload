package ds

import "sync"

type ArrayList[T any] struct {
	mu   sync.RWMutex
	data []T
	eq   func(a, b T) bool
}

func NewArrayList[T any](eq func(a, b T) bool) *ArrayList[T] {
	return &ArrayList[T]{
		data: make([]T, 0),
		eq:   eq,
	}
}

func (l *ArrayList[T]) Add(item T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data = append(l.data, item)
}

func (l *ArrayList[T]) AddAt(index int, item T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index > len(l.data) {
		return false
	}

	l.data = append(l.data[:index], append([]T{item}, l.data[index:]...)...)
	return true
}

func (l *ArrayList[T]) Get(index int) (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if index < 0 || index >= len(l.data) {
		var zero T
		return zero, false
	}

	return l.data[index], true
}

func (l *ArrayList[T]) Set(index int, item T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index >= len(l.data) {
		return false
	}

	l.data[index] = item
	return true
}

func (l *ArrayList[T]) RemoveAt(index int) (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index >= len(l.data) {
		var zero T
		return zero, false
	}

	val := l.data[index]
	l.data = append(l.data[:index], l.data[index+1:]...)
	return val, true
}

func (l *ArrayList[T]) Remove(item T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, v := range l.data {
		if l.eq(v, item) {
			l.data = append(l.data[:i], l.data[i+1:]...)
			return true
		}
	}

	return false
}

func (l *ArrayList[T]) Contains(item T) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, v := range l.data {
		if l.eq(v, item) {
			return true
		}
	}

	return false
}

func (l *ArrayList[T]) IndexOf(item T) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for i, v := range l.data {
		if l.eq(v, item) {
			return i
		}
	}

	return -1
}

func (l *ArrayList[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.data = make([]T, 0)
}

func (l *ArrayList[T]) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.data)
}

func (l *ArrayList[T]) IsEmpty() bool {
	return l.Size() == 0
}

func (l *ArrayList[T]) ToSlice() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]T, len(l.data))
	copy(out, l.data)

	return out
}
