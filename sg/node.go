package sg

import (
    "fmt"
)

type Node interface {
    GetChildren() [] Node
}

type TransformNode struct {
    Dx float32
    Dy float32
    Scale float32
    Rotation float32

    Children []Node
}

type RectangleNode struct {
    X float32
    Y float32
    W float32
    H float32

    R float32
    G float32
    B float32
    A float32

    Children []Node
}

type TextureNode struct {
    X float32
    Y float32
    W float32
    H float32

    TX float32
    TY float32
    TW float32
    TH float32

    Texture uint;

    Children []Node
}



func (this RectangleNode) String() string {
    return fmt.Sprintf("Rectangle(%.1f,%.1f %.1fx%.1f - rgba=%.1f %.1f %.1f %.1f)",
            this.X,
            this.Y,
            this.W,
            this.H,
            this.R,
            this.G,
            this.B,
            this.A);
}

func (this *RectangleNode) SetColor(r , g, b, a float32) {
    this.R = r
    this.G = g
    this.B = b
    this.A = a
}

func (this *RectangleNode) SetGeometry(x, y, w, h float32) {
    this.X = x
    this.Y = y
    this.W = w
    this.H = h
}

func (this RectangleNode) GetChildren() []Node {
    return this.Children
}


func (this TransformNode) String() string {
    return fmt.Sprintf("Transform(d=%.1f,%.1f, s=%.1f, r=%.2f)",
                       this.Dx, this.Dy,
                       this.Scale, this.Rotation);
}

func (this TransformNode) GetChildren() []Node {
    return this.Children
}




