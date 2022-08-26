package linearquadtree

import (
	"sync"
)

type CollisionPair struct {
	Obj1 *TreeObject
	Obj2 *TreeObject
}

type Cell struct {
	Latest *TreeObject
}

type TreeObject struct {
	Cell   *Cell
	Object interface{}
	Before *TreeObject
	Next   *TreeObject
}

const TreeMaxLevel = 10

type CLiner4TreeManager struct {
	Cells      []*Cell
	Pow        [TreeMaxLevel + 1]int32
	Width      float64 // 領域のX軸幅
	Height     float64 // 領域のY軸幅
	Left       float64 // 領域の左側（X軸最小値）
	Top        float64 // 領域の上側（Y軸最小値）
	UnitWidth  float64 // 最小レベル空間の幅単位
	UnitHeight float64 // 最小レベル空間の高単位
	CellNum    int32   // 空間数
	Level      int32
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

func (m *CLiner4TreeManager) Init(left, top, right, bottom float64) {

	level := int32(9)
	// 各レベルでの空間数を算出
	m.Pow[0] = 1
	for i := 1; i < TreeMaxLevel+1; i++ {
		m.Pow[i] = m.Pow[i-1] * 4
	}

	m.CellNum = (m.Pow[level+1] - 1) / 3
	m.Cells = make([]*Cell, m.CellNum)

	// 有効領域を登録
	m.Left = left
	m.Top = top
	m.Width = right - left
	m.Height = bottom - top
	m.UnitWidth = m.Width / float64((int(1) << level))
	m.UnitHeight = m.Height / float64((int(1) << level))
	m.Level = level
}

func (m *CLiner4TreeManager) Register(left, top, right, bottom float64, obj *TreeObject) {
	// オブジェクトの境界範囲から登録モートン番号を算出
	mortonNo := m.getMortonNumber(left, top, right, bottom)
	if mortonNo < m.CellNum {
		// 空間が無い場合は新規作成
		if m.Cells[mortonNo] == nil {
			m.createNewCell(mortonNo)
		}
		m.Cells[mortonNo].Push(obj)
	}
}

func (m *CLiner4TreeManager) GetCollisionList() []CollisionPair {
	pairs := make([]CollisionPair, 0)
	if m.Cells[0] == nil {
		return pairs
	}

	var stack *TreeObjectStack = NewStack(1000)

	m._getCollisionList(0, pairs, stack)

	return pairs
}

func (m *CLiner4TreeManager) _getCollisionList(elem int32, pairs []CollisionPair, stack *TreeObjectStack) {
	obj1 := m.Cells[elem].Latest

	for obj1 != nil {
		obj2 := obj1.Next
		// 空間内の衝突可能性リスト
		for obj2 != nil {
			pairs = append(pairs, CollisionPair{obj1, obj2})
			obj2 = obj2.Next
		}

		// スタックとの衝突可能性リスト
		for _, stackObj := range stack.data {
			pairs = append(pairs, CollisionPair{obj1, stackObj})
		}

		obj1 = obj1.Next
	}

	childFlag := false
	objNum := 0
	var nextElem int32
	for i := 0; i < 4; i++ {
		nextElem = elem*4 + 1 + int32(i)
		if nextElem < m.CellNum && m.Cells[nextElem] != nil {
			if !childFlag {
				obj1 := m.Cells[elem].Latest

				for obj1 != nil {
					stack.Push(obj1)
					objNum++
					obj1 = obj1.Next
				}
				childFlag = true
				m._getCollisionList(nextElem, pairs, stack)
			}
		}
	}

	// stackからオブジェクトを外す
	if childFlag {
		for i := 0; i < objNum; i++ {
			stack.Pop()
		}
	}
}

func (m *CLiner4TreeManager) createNewCell(cellNum int32) {
	for m.Cells[cellNum] == nil {
		m.Cells[cellNum] = &Cell{}
		cellNum = (cellNum - 1) >> 2
		if cellNum >= m.CellNum {
			break
		}
	}
}

// 座標→線形4分木要素番号変換関数
func (m *CLiner4TreeManager) getPointElem(pos_x, pos_y float64) int32 {
	return Get2DMortonNumber((int32)((pos_x-m.Left)/m.UnitWidth), (int32)((pos_y-m.Top)/m.UnitHeight))
}

func (m *CLiner4TreeManager) getMortonNumber(left, top, right, bottom float64) int32 {
	// 最小レベルにおける各軸位置を算出
	leftTop := m.getPointElem(left, top)
	rightBottom := m.getPointElem(right, bottom)

	// 空間番号の排他的論理和から
	// 所属レベルを算出
	def := rightBottom ^ leftTop
	hiLevel := 0
	for i := 0; i < int(m.Level); i++ {
		check := (def >> (i * 2)) & 0x3
		if check != 0 {
			hiLevel = i + 1
		}
	}
	spaceNum := rightBottom >> (int32(hiLevel) * 2)
	addNum := m.Pow[m.Level-int32(hiLevel)-1] / 3
	spaceNum += addNum
	if spaceNum > m.CellNum {
		panic("over cell no")
	}
	return spaceNum
}
