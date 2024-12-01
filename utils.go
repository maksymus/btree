package main

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		panic("empty stack")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack[T]) Top() T {
	if len(s.items) == 0 {
		panic("empty stack")
	}

	return s.items[len(s.items)-1]
}

func (s *Stack[T]) Empty() bool {
	return len(s.items) == 0
}
