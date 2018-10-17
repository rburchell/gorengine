package main

import (
	"fmt"
	"github.com/MichaelTJones/pcg"
	"github.com/sletta/gorengine/sg"
	"log"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

func init() {
	runtime.LockOSThread()
}

/*
func dumpTree(node sg.Node, level int) {

	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
	fmt.Println(" -", node)

	for _, child := range node.GetChildren() {
		dumpTree(child, level+1)
	}
}
*/

func build(nodeChan chan sg.Node) {
	nodeChan <- sg.TransformNode{
		Scale: 1,
	}

	nodeChan <- sg.RectangleNode{
		X: 100,
		Y: 100,
		W: 200,
		H: 100,
		R: 1,
		G: 0,
		B: 1,
		A: 1,
	}

	var wg sync.WaitGroup
	const childSenders = 10
	const childBatch = 500
	const childCount = childSenders * childBatch
	wg.Add(childSenders)

	for i := 0; i < childSenders; i++ {
		go func() {
			var pcgrand = pcg.NewPCG32()
			pcgrand.Seed(uint64(rand.Int63()), uint64(rand.Int63()))

			for i := 0; i < childBatch; i++ {
				nodeChan <- sg.RectangleNode{
					X: float32(float64(pcgrand.Random())/float64(math.MaxUint32)) * 800,
					Y: float32(float64(pcgrand.Random())/float64(math.MaxUint32)) * 410,
					W: float32(float64(pcgrand.Random())/float64(math.MaxUint32)) * 200,
					H: float32(float64(pcgrand.Random())/float64(math.MaxUint32)) * 200,
					R: float32(float64(pcgrand.Random()) / float64(math.MaxUint32)),
					G: float32(float64(pcgrand.Random()) / float64(math.MaxUint32)),
					B: float32(float64(pcgrand.Random()) / float64(math.MaxUint32)),
					A: float32(float64(pcgrand.Random()) / float64(math.MaxUint32)),
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
	close(nodeChan)

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
			go build(renderer.RenderChan)
			//dumpTree(root, 0)
			renderer.SetClearColor(1, 1, 1, 1)
			renderer.Render()
		}
		renderer.Destroy()
	} else {
		fmt.Println("Failed to create renderer... aborting main()")
	}
}
