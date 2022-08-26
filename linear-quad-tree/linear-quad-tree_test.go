package linearquadtree

import "testing"

func TestGet2DMortonNumber(t *testing.T) {
	number := Get2DMortonNumber(3, 6)
	println(number)
	if number != 45 {
		t.Error(`no`)
	}

	number = Get2DMortonNumber(6, 5)
	println(number)
	if number != 54 {
		t.Error(`no`)
	}
}

func TestGetPointElem(t *testing.T) {
	manager := Liner4TreeManager{}
	manager.Init(3, 0, 0, 100, 100) // 空間は64個

	point := manager.getPointElem(51, 51)
	if point != 48 {
		t.Error()
	}

	point = manager.getPointElem(0, 0)
	if point != 0 {
		t.Error()
	}

	point = manager.getPointElem(99, 99)
	if point != 63 {
		t.Error()
	}

	point = manager.getPointElem(100, 100) // 0からlength100なので、100は範囲外
	if point <= 63 {
		t.Error()
	}
}

func TestGetMortonNumber(t *testing.T) {
	manager := Liner4TreeManager{}

	// level0 1個
	// level1 4個
	// level2 16個
	// level3 64個
	manager.Init(3, 0, 0, 100, 100)

	// level0 0に所属するはず
	point := manager.getMortonNumber(0, 0, 99, 99)
	println(point)
	if point != 0 {
		t.Error()
	}

	// level1 1に所属するはず
	point = manager.getMortonNumber(51, 0, 99, 48)
	println(point)
	// level0(1) + level1(1)
	if point != 2 {
		t.Error()
	}

	// level2 5に所属するはず
	// level0(1) + level1(4) + level2(5)
	point = manager.getMortonNumber(76, 0, 99, 24)
	println(point)
	if point != 10 {
		t.Error()
	}

	// level3 63に所属するはず
	// level0(1)+level1(4)+level2(16)+level3(63) = 84
	point = manager.getMortonNumber(98, 98, 99, 99)
	println(point)
	if point != 84 {
		t.Error()
	}
}
