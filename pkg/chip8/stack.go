package chip8

import "fmt"

type stack struct {
	Stack []uint16
	Top   int
}

func newStack(size int) stack {
	return stack{
		Stack: make([]uint16, size),
		Top:   0,
	}
}

func (s *stack) Push(value uint16) error {
	if s.Top == len(s.Stack) {
		return fmt.Errorf("stack overflow: could not push value %d to stack as limit %d is already reached", value, len(s.Stack))
	}

	s.Top++
	s.Stack[s.Top] = value

	return nil
}

func (s *stack) Pop() (uint16, error) {
	if s.Top == 0 {
		return 0, fmt.Errorf("stack underflow: could not pop value from stack as bottom is already reached")
	}

	value := s.Stack[s.Top]
	s.Top--

	return value, nil
}
