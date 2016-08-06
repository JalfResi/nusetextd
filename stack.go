package main

import "sync"

// Stack struct
type Stack struct {
	sync.Mutex
	stack []Worker
}

// Inc method
func (s *Stack) Inc(count int) []Worker {
	temp := make([]Worker, count)

	for i := range temp {
		temp[i] = make(Worker, 1)
	}

	s.Lock()
	defer s.Unlock()
	s.stack = append(s.stack, temp...)

	return temp
}

// Dec method
func (s *Stack) Dec(count int) []Worker {
	s.Lock()
	defer s.Unlock()

	pos := len(s.stack) - count

	n := make([]Worker, count)
	copy(n, s.stack[pos:])

	for i := range s.stack[pos:] {
		s.stack[i] = nil
	}
	s.stack = s.stack[:pos]

	return n
}

// Len method
func (s *Stack) Len() int {
	s.Lock()
	defer s.Unlock()

	return len(s.stack)
}
