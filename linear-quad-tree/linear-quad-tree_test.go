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
