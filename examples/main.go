package main

import (
    "fmt"
    "runtime"
    "github.com/sletta/gorengine/sg"
)

func init() {
    runtime.LockOSThread()
}

func dumpTree(node sg.Node, level int ) {

    for i := 0; i<level; i++ {
        fmt.Print(" ")
    }
    fmt.Println(" -", node)

    for _, child := range node.GetChildren() {
        dumpTree(child, level + 1)
    }
}

func build() sg.Node {
    return sg.TransformNode{
            Scale: 1,
            Children: [] sg.Node{
                sg.RectangleNode{
                    X: 100,
                    Y: 100,
                    W: 200,
                    H: 100,
                    R: 0,
                    G: 0,
                    B: 1,
                    A: 1,
                    Children: [] sg.Node{
                        sg.RectangleNode{
                            X: 110,
                            Y: 110,
                            W: 20,
                            H: 10,
                            R: 1,
                            G: 1,
                            B: 0,
                            A: 1,
                        },
                    },
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

    var root sg.Node = build()
    dumpTree(root, 0)

    var renderer *sg.Renderer = sg.CreateRenderer()

    if renderer != nil {
        for !renderer.ShouldClose() {
            renderer.SetClearColor(1, 1, 1, 1)
            renderer.Render(root)
        }
        renderer.Destroy()
    } else {
        fmt.Println("Failed to create renderer... aborting main()");
    }
}
