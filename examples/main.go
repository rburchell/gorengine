package main

import (
	"fmt"
	"github.com/sletta/gorengine/sg"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

type randStorage struct {
	randBuf     []float32
	currentRand int
}

func (this *randStorage) shuffle() {
	if len(this.randBuf) == 0 {
		// init on first acquire
		for i := 0; i < 5000; i++ {
			this.randBuf = append(this.randBuf, rand.Float32())
		}
	}
	this.currentRand = rand.Intn(len(this.randBuf))
}

func (this *randStorage) acquire() float32 {
	this.currentRand += 1
	if this.currentRand == len(this.randBuf) {
		this.currentRand = 0
	}
	return this.randBuf[this.currentRand]
}

var randCache randStorage

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
	randCache.shuffle()
	const childCount = 10000
	childs := make([]sg.Node, childCount)
	for i := 0; i < childCount; i++ {
		childs[i] = sg.RectangleNode{
			X: randCache.acquire() * 800,
			Y: randCache.acquire() * 410,
			W: randCache.acquire() * 200,
			H: randCache.acquire() * 200,
			R: randCache.acquire(),
			G: randCache.acquire(),
			B: randCache.acquire(),
			A: randCache.acquire(),
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
