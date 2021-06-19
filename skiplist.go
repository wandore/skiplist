package skiplist

import "math/rand"

const (
	maxLevel int     = 16   // 跳表最大深度
	p        float32 = 0.25 // 概率因子
)

// 生成随机深度
func randomLevel() int {
	level := 1
	for rand.Float32() < p && level < maxLevel {
		level++
	}
	return level
}

// 跳表节点
type Node struct {
	Score   float64     // 节点权值
	Value   interface{} // 节点值
	forward []*Node     // 指向节点列表
}

// 构造节点
func newElement(score float64, value interface{}, level int) *Node {
	return &Node{
		Score:   score,
		Value:   value,
		forward: make([]*Node, level),
	}
}

// 获取指向的第一个节点
func (e *Node) Next() *Node {
	if e != nil {
		return e.forward[0]
	}
	return nil
}

// 跳表
type SkipList struct {
	header *Node // 虚拟头节点 dummy，不计入长度/深度
	len    int   // 跳表长度
	level  int   // 跳表深度
}

// 构造跳表
func New() *SkipList {
	return &SkipList{header: &Node{forward: make([]*Node, maxLevel)}}
}

// 获取跳表长度
func (s *SkipList) Size() int {
	return s.len
}

// 获取跳表第一个节点
func (s *SkipList) Front() *Node {
	return s.header.forward[0]
}

// 根据权值查找节点
func (s *SkipList) Search(score float64) (*Node, bool) {
	// 从 dummy 节点开始遍历
	e := s.header
	// 从高往低查找
	for i := s.level - 1; i >= 0; i-- {
		// 在当前层获取到最后一个小于权值的节点，即查询权值可能的前置节点
		for e.forward[i] != nil && e.forward[i].Score < score {
			e = e.forward[i]
		}
	}
	// 在最底层，即第0层获取到前置节点的后一节点
	e = e.forward[0]

	// 权值存在
	if e != nil && e.Score == score {
		return e, true
	}

	// 权值不存在
	return nil, false
}

// 插入节点
func (s *SkipList) Insert(score float64, value interface{}) *Node {
	update := make([]*Node, maxLevel)

	// 获取插入位置
	e := s.header
	for i := s.level - 1; i >= 0; i-- {
		for e.forward[i] != nil && e.forward[i].Score < score {
			e = e.forward[i]
		}
		// 将每一层插入节点的前置节点保存到 update 数组
		update[i] = e
	}
	// 最底层，即第0层
	e = e.forward[0]

	// 权值已经存在
	if e != nil && e.Score == score {
		e.Value = value
		return e
	}

	// 确定插入节点的深度，min(maxLevel, s.level + 1)，每次最多增加一层
	level := randomLevel()
	if level > s.level {
		level = s.level + 1
		update[s.level] = s.header
		s.level = level
	}

	// 构造节点
	node := newElement(score, value, level)
	
	// 将节点插入跳表，从低往高每一层的前置节点->node->前置节点的后一节点
	for i := 0; i < level; i++ {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}

	// 更新跳表长度
	s.len++

	return node
}

// 删除节点
func (s *SkipList) Delete(score float64) *Node {
	update := make([]*Node, maxLevel)

	// 获取插入位置
	e := s.header
	for i := s.level - 1; i >= 0; i-- {
		for e.forward[i] != nil && e.forward[i].Score < score {
			e = e.forward[i]
		}
		// 将每一层插入节点的前置节点保存到 update 数组
		update[i] = e
	}
	// 最底层，即第0层
	e = e.forward[0]

	// 权值不存在
	if e == nil || e.Score != score {
		return nil
	}

	// 从低往高
	for i := 0; i < s.level; i++ {
		// 当前层的后置节点不为指定权值，证明当前高度大于权值节点的深度，跳出循环
		if update[i].forward[i] != e {
			break
		}
		update[i].forward[i] = e.forward[i]
	}

	// 更新跳表长度
	s.len--

	return e
}
