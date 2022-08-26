package linearquadtree

import (
	"fmt"
	"sync"
)

type CollisionPair struct {
	Obj1 *TreeObject
	Obj2 *TreeObject
}

type Cell struct {
	Latest *TreeObject
	Mu     sync.Mutex
}

type TreeObject struct {
	Cell   *Cell
	Object interface{}
	Before *TreeObject
	Next   *TreeObject
}

const TreeMaxLevel = 10

type Liner4TreeManager struct {
	cells      []*Cell
	pow        [TreeMaxLevel + 1]int32
	width      float64 // 領域のX軸幅
	height     float64 // 領域のY軸幅
	left       float64 // 領域の左側（X軸最小値）
	top        float64 // 領域の上側（Y軸最小値）
	unitWidth  float64 // 最小レベル空間の幅単位
	unitHeight float64 // 最小レベル空間の高単位
	cellNum    int32   // 空間数
	level      int32
	Mu         sync.Mutex
}

func BitSeparate32(n int32) int32 {
	n = (n | (n << 8)) & 0x00ff00ff
	n = (n | (n << 4)) & 0x0f0f0f0f
	n = (n | (n << 2)) & 0x33333333
	return (n | (n << 1)) & 0x55555555
}

func Get2DMortonNumber(x, y int32) int32 {
	return (BitSeparate32(x) | (BitSeparate32(y) << 1))
}

func (c *Cell) Push(obj *TreeObject) {
	if obj == nil {
		return // 無効オブジェクトは登録しない
	}

	if obj.Cell == c {
		return // 2重登録チェック
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.Latest == nil {
		c.Latest = obj
	} else {
		obj.Next = c.Latest
		c.Latest.Before = obj
		c.Latest = obj
	}

	obj.Cell = c
}

func (c *Cell) OnRemove(obj *TreeObject) {
	if c.Latest == obj {
		c.Latest = obj.Next
	}
}

func (obj *TreeObject) Remove() {
	// すでに逸脱している時は処理終了
	if obj.Cell == nil {
		return
	}

	obj.Cell.Mu.Lock()
	defer obj.Cell.Mu.Unlock()

	// 自分を登録している空間に自身を通知
	obj.Cell.OnRemove(obj)

	// 逸脱処理
	// 前後のオブジェクトを結びつける
	if obj.Before != nil {
		obj.Before.Next = obj.Next
	}

	if obj.Next != nil {
		obj.Next.Before = obj.Before
	}

	obj.Before = nil
	obj.Next = nil
	obj.Cell = nil
}

func (m *Liner4TreeManager) Init(level int32, left, top, right, bottom float64) {
	// 各レベルでの空間数を算出
	m.pow[0] = 1
	for i := 1; i < TreeMaxLevel+1; i++ {
		m.pow[i] = m.pow[i-1] * 4
	}

	m.cellNum = (m.pow[level+1] - 1) / 3
	m.cells = make([]*Cell, m.cellNum)

	// 有効領域を登録
	m.left = left
	m.top = top
	m.width = right - left
	m.height = bottom - top
	m.unitWidth = m.width / float64((int(1) << level))
	m.unitHeight = m.height / float64((int(1) << level))
	m.level = level
}

func (m *Liner4TreeManager) Register(left, top, right, bottom float64, treeObj *TreeObject) error {
	// オブジェクトの境界範囲から登録モートン番号を算出
	mortonNo := m.getMortonNumber(left, top, right, bottom)
	if mortonNo < m.cellNum {
		// 空間が無い場合は新規作成
		if m.cells[mortonNo] == nil {
			m.createNewCell(mortonNo)
		}
		m.cells[mortonNo].Push(treeObj)
		return nil
	} else {
		return fmt.Errorf("object over range cellnum")
	}
}

func (m *Liner4TreeManager) GetCollisionList() []CollisionPair {
	pairs := make([]CollisionPair, 0, 1000000)
	if m.cells[0] == nil {
		return pairs
	}

	var stack *TreeObjectStack = NewStack(1000000)

	m._getCollisionList(0, &pairs, stack)

	return pairs
}

func (m *Liner4TreeManager) _getCollisionList(elem int32, pairs *[]CollisionPair, stack *TreeObjectStack) {
	obj1 := m.cells[elem].Latest
	for obj1 != nil {
		obj2 := obj1.Next
		// 空間内の衝突可能性リスト
		for obj2 != nil {
			*pairs = append(*pairs, CollisionPair{obj1, obj2})
			obj2 = obj2.Next
		}

		// スタックとの衝突可能性リスト
		for _, stackObj := range stack.data {
			*pairs = append(*pairs, CollisionPair{obj1, stackObj})
		}

		obj1 = obj1.Next
	}

	childFlag := false
	objNum := 0
	var nextElem int32
	for i := 0; i < 4; i++ {
		nextElem = elem*4 + 1 + int32(i)
		if nextElem < m.cellNum && m.cells[nextElem] != nil {
			if !childFlag {
				obj1 := m.cells[elem].Latest

				for obj1 != nil {
					stack.Push(obj1)
					objNum++
					obj1 = obj1.Next
				}
			}
			childFlag = true
			m._getCollisionList(nextElem, pairs, stack)
		}
	}

	// stackからオブジェクトを外す
	if childFlag {
		for i := 0; i < objNum; i++ {
			stack.Pop()
		}
	}
}

func (m *Liner4TreeManager) createNewCell(cellNum int32) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	for m.cells[cellNum] == nil {
		m.cells[cellNum] = &Cell{}
		cellNum = (cellNum - 1) >> 2
		if cellNum < 0 {
			break
		}
	}
}

// 座標→線形4分木要素番号変換関数
func (m *Liner4TreeManager) getPointElem(x, y float64) int32 {
	return Get2DMortonNumber((int32)((x-m.left)/m.unitWidth), (int32)((y-m.top)/m.unitHeight))
}

func (m *Liner4TreeManager) getMortonNumber(left, top, right, bottom float64) int32 {
	// 最小レベルにおける各軸位置を算出
	leftTop := m.getPointElem(left, top)
	rightBottom := m.getPointElem(right, bottom)

	// 空間番号の排他的論理和から
	// 所属レベルを算出
	def := rightBottom ^ leftTop
	hiLevel := 0
	for i := 0; i < int(m.level); i++ {
		check := (def >> (i * 2)) & 0x3
		if check != 0 {
			hiLevel = i + 1
		}
	}
	spaceNum := rightBottom >> (int32(hiLevel) * 2)
	addNum := (m.pow[m.level-int32(hiLevel)] - 1) / 3
	spaceNum += addNum
	return spaceNum
}
