package main

import (
	linearquadtree "github.com/hogesako/go-linear-quad-tree/linear-quad-tree"
)

type Circle struct {
	Center int
}

func main() {
	manager := linearquadtree.CLiner4TreeManager{}

	manager.Init(0, 0, 100, 100)
	obj := linearquadtree.TreeObject{}
	obj.Object = Circle{1}
	manager.Register(1, 1, 3, 3, &obj)

	obj = linearquadtree.TreeObject{}
	obj.Object = Circle{2}
	manager.Register(1, 1, 3, 3, &obj)

	obj = linearquadtree.TreeObject{}
	obj.Object = Circle{3}
	manager.Register(50, 50, 60, 60, &obj)

	list := manager.GetCollisionList()
	println("count", len(list))
	for _, v := range list {
		println(v.Obj1.Object.(Circle).Center, v.Obj2.Object.(Circle).Center)
	}
}
