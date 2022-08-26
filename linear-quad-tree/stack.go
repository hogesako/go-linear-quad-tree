package linearquadtree

import "fmt"

type TreeObjectStack struct {
	data []*TreeObject
	size int
}

func NewStack(cap int) *TreeObjectStack {
	return &TreeObjectStack{data: make([]*TreeObject, 0, cap), size: 0}
}

// Push adds a new element at the end of the stack
func (s *TreeObjectStack) Push(n *TreeObject) {
	s.data = append(s.data, n)
	s.size++
}

// Pop removes the last element from stack
func (s *TreeObjectStack) Pop() bool {
	if s.IsEmpty() {
		return false
	}
	s.size--
	s.data = s.data[:s.size]
	return true
}

// Top returns the last element of stack
func (s *TreeObjectStack) Top() *TreeObject {
	return s.data[s.size-1]
}

// IsEmpty checks if the stack is empty
func (s *TreeObjectStack) IsEmpty() bool {
	return s.size == 0
}

// String implements Stringer interface
func (s *TreeObjectStack) String() string {
	return fmt.Sprint(s.data)
}
