package skiplist

import "math/rand"

const (
	maxLevel int     = 16
	p        float32 = 0.25
)

func randomLevel() int {
	level := 1
	for rand.Float32() < p && level < maxLevel {
		level++
	}
	return level
}

type Node struct {
	Score   float64
	Value   interface{}
	forward []*Node
}

func newElement(score float64, value interface{}, level int) *Node {
	return &Node{
		Score:   score,
		Value:   value,
		forward: make([]*Node, level),
	}
}

func (e *Node) Next() *Node {
	if e != nil {
		return e.forward[0]
	}
	return nil
}

type SkipList struct {
	header *Node
	len    int
	level  int
}

func New() *SkipList {
	return &SkipList{header: &Node{forward: make([]*Node, maxLevel)}}
}

func (s *SkipList) Size() int {
	return s.len
}

func (s *SkipList) Front() *Node {
	return s.header.forward[0]
}

func (s *SkipList) Search(score float64) (*Node, bool) {
	e := s.header
	for i := s.level - 1; i >= 0; i-- {
		for e.forward[i] != nil && e.forward[i].Score < score {
			e = e.forward[i]
		}
	}
	e = e.forward[0]
	if e != nil && e.Score == score {
		return e, true
	}
	return nil, false
}
