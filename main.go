package main

import (
	linearquadtree "github.com/hogesako/go-linear-quad-tree/linear-quad-tree"
)

func main() {
	manager := linearquadtree.CLiner4TreeManager{}

	manager.Init(0, 0, 100, 100)
}
