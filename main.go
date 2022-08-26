package main

import (
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	linearquadtree "github.com/hogesako/go-linear-quad-tree/linear-quad-tree"
)

type Circle struct {
	Center int
}

func main() {
	manager := linearquadtree.Liner4TreeManager{}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	manager.Init(9, 0, 0, 100, 100)
	obj := linearquadtree.TreeObject{}
	obj.Object = Circle{1}
	manager.Register(1, 1, 3, 3, &obj)

	obj2 := linearquadtree.TreeObject{}
	obj2.Object = Circle{2}
	manager.Register(1, 1, 3, 3, &obj2)

	obj3 := linearquadtree.TreeObject{}
	obj3.Object = Circle{3}
	manager.Register(51, 51, 80, 80, &obj3)

	now := time.Now()
	// for i := 0; i < 1000; i++ {
	// 	rand.Seed(time.Now().UnixNano())

	// 	obj4 := linearquadtree.TreeObject{}
	// 	obj4.Object = Circle{i}
	// 	obj4.Next = nil

	// 	point := rand.Float64() * 99
	// 	manager.Register(point, point, point, point, &obj4)
	// }
	println(time.Since(now).String())

	list := manager.GetAllCollisionList()
	println("count", len(list))

}
