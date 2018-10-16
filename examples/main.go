package main

import (
	"fmt"
	"github.com/sletta/gorengine/sg"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

func dumpTree(node sg.Node, level int) {

	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
	fmt.Println(" -", node)

	for _, child := range node.GetChildren() {
		dumpTree(child, level+1)
	}
}

func build() sg.Node {
	const childCount = 10000
	childs := make([]sg.Node, childCount)
	for i := 0; i < childCount; i++ {
		childs[i] = sg.RectangleNode{
			X: rand.Float32() * 800,
			Y: rand.Float32() * 410,
			W: rand.Float32() * 200,
			H: rand.Float32() * 200,
			R: float32(rand.NormFloat64()),
			G: float32(rand.NormFloat64()),
			B: float32(rand.NormFloat64()),
			A: float32(rand.NormFloat64()),
		}
	}

	return sg.TransformNode{
		Scale: 1,
		Children: []sg.Node{
			sg.RectangleNode{
				X:        100,
				Y:        100,
				W:        200,
				H:        100,
				R:        1,
				G:        0,
				B:        1,
				A:        1,
				Children: childs,
			},
		},
	}
	// var root sg.TransformNode = sg.TransformNode{ Scale: 1 }
	// var rectangle sg.RectangleNode = sg.RectangleNode{}
	// rectangle.SetGeometry(100, 200, 500, 300)
	// rectangle.SetColor(1, 0, 0.5, 1)

	// var childRect sg.RectangleNode = sg.RectangleNode{}
	// childRect.SetGeometry(10, 10, 10, 10)
	// childRect.SetColor(0, 1, 0, 1)

	// root.AddChild(rectangle)
	// rectangle.AddChild(childRect)

	// childRect.SetGeometry(20, 20, 20, 20)

	// return root
}

func main() {

	var renderer *sg.Renderer = sg.CreateRenderer()

	if renderer != nil {
		lastRender := time.Now()
		for !renderer.ShouldClose() {
			log.Printf("Time since last frame: %s", time.Since(lastRender))
			lastRender = time.Now()
			var root sg.Node = build()
			//dumpTree(root, 0)
			renderer.SetClearColor(1, 1, 1, 1)
			renderer.Render(root)
		}
		renderer.Destroy()
	} else {
		fmt.Println("Failed to create renderer... aborting main()")
	}
}
