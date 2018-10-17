package sg

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"strings"
)

type Renderer struct {
	program      uint32
	vertexBuffer uint32
	indexBuffer  uint32

	vertexData []float32
	indexData  []uint16

	window *glfw.Window

	RenderChan chan Node
}

func CreateRenderer() *Renderer {
	fmt.Println("Creating renderer")
	if err := glfw.Init(); err != nil {
		fmt.Println("Failed to initialize GLFW...")
		return nil
	}

	var windowWidth = 800
	var windowHeight = 600

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var renderer Renderer = Renderer{}
	renderer.RenderChan = make(chan Node)
	var err error

	renderer.window, err = glfw.CreateWindow(windowWidth, windowHeight, "gorengine", nil, nil)
	if err != nil {
		panic(err)
	}
	renderer.window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		fmt.Println("Failed to initialize GL 2.1...")
		return nil
	}

	renderer.program, err = newProgram(vertexSourceCode, fragmentSourceCode)
	if err != nil {
		fmt.Println("Failed to create shader program", err)
		return nil
	}
	gl.UseProgram(renderer.program)
	var uniformWindowSize = gl.GetUniformLocation(renderer.program, gl.Str("WindowSize\x00"))
	gl.Uniform2f(int32(uniformWindowSize), float32(windowWidth), float32(windowHeight))
	gl.UseProgram(0)

	gl.GenBuffers(1, &renderer.vertexBuffer)
	gl.GenBuffers(1, &renderer.indexBuffer)

	fmt.Println("Renderer created ok..")

	return &renderer
}

func (this *Renderer) Destroy() {
	// ### Delete program and buffers...
	glfw.Terminate()
}

func (this *Renderer) ShouldClose() bool {
	return this.window.ShouldClose()
}

func (this *Renderer) renderBuildBuffers(nodeChan chan Node) int {
	this.vertexData = this.vertexData[0:0]
	this.indexData = this.indexData[0:0]

	rectCount := 0

	for node := range nodeChan {
		if rn, ok := node.(RectangleNode); ok {
			rectCount++
			var firstIndex uint16 = uint16(len(this.vertexData) / 9)
			this.indexData = append(this.indexData, firstIndex+0)
			this.indexData = append(this.indexData, firstIndex+1)
			this.indexData = append(this.indexData, firstIndex+2)
			this.indexData = append(this.indexData, firstIndex+1)
			this.indexData = append(this.indexData, firstIndex+3)
			this.indexData = append(this.indexData, firstIndex+2)

			this.vertexData = append(this.vertexData, rn.X)
			this.vertexData = append(this.vertexData, rn.Y)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, rn.R)
			this.vertexData = append(this.vertexData, rn.G)
			this.vertexData = append(this.vertexData, rn.B)
			this.vertexData = append(this.vertexData, rn.A)
			this.vertexData = append(this.vertexData, 0.0)

			this.vertexData = append(this.vertexData, rn.X+rn.W)
			this.vertexData = append(this.vertexData, rn.Y)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, rn.R)
			this.vertexData = append(this.vertexData, rn.G)
			this.vertexData = append(this.vertexData, rn.B)
			this.vertexData = append(this.vertexData, rn.A)
			this.vertexData = append(this.vertexData, 0.0)

			this.vertexData = append(this.vertexData, rn.X)
			this.vertexData = append(this.vertexData, rn.Y+rn.H)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, rn.R)
			this.vertexData = append(this.vertexData, rn.G)
			this.vertexData = append(this.vertexData, rn.B)
			this.vertexData = append(this.vertexData, rn.A)
			this.vertexData = append(this.vertexData, 0.0)

			this.vertexData = append(this.vertexData, rn.X+rn.W)
			this.vertexData = append(this.vertexData, rn.Y+rn.H)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, 0)
			this.vertexData = append(this.vertexData, rn.R)
			this.vertexData = append(this.vertexData, rn.G)
			this.vertexData = append(this.vertexData, rn.B)
			this.vertexData = append(this.vertexData, rn.A)
			this.vertexData = append(this.vertexData, 0.0)
		}
	}

	return rectCount
}

func (this *Renderer) Render() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	count := this.renderBuildBuffers(this.RenderChan)

	// fmt.Println("Render:")
	// for i:=0; i<int(vertexCount); i++ {
	//     fmt.Print(" - vertex: ", i, ": ");
	//     for v:=0; v<9; v++ {
	//         fmt.Print(this.vertexData[i * 9 + v], ", ")
	//     }
	//     fmt.Println()
	// }
	// for i:=0; i<int(count); i++ {
	//     fmt.Printf(" - index %d: %d %d %d %d %d %d\n", i,
	//                 this.indexData[i * 6 + 0],
	//                 this.indexData[i * 6 + 1],
	//                 this.indexData[i * 6 + 2],
	//                 this.indexData[i * 6 + 3],
	//                 this.indexData[i * 6 + 4],
	//                 this.indexData[i * 6 + 5])
	// }

	vertexCount := count * 4
	indexCount := count * 6
	vertexBufferSize := vertexCount * 9
	indexBufferSize := indexCount

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, int(vertexBufferSize*4), gl.Ptr(this.vertexData), gl.STREAM_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indexBufferSize*2), gl.Ptr(this.indexData), gl.STREAM_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1) // texture coords
	gl.EnableVertexAttribArray(2) // color
	gl.EnableVertexAttribArray(3) // type..

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 9*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 9*4, gl.PtrOffset(2*4))
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(4*4))
	gl.VertexAttribPointer(3, 1, gl.FLOAT, false, 9*4, gl.PtrOffset(8*4))

	gl.UseProgram(this.program)
	gl.DrawElements(gl.TRIANGLES, int32(indexBufferSize), gl.UNSIGNED_SHORT, gl.PtrOffset(0))

	gl.UseProgram(0)

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)

	this.window.SwapBuffers()
	glfw.PollEvents()

	this.RenderChan = make(chan Node)
}

func (this *Renderer) SetClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexSourceCode = `
attribute vec2 Vertex;
attribute vec2 TexCoord;
attribute vec4 Color;
attribute float Type;

uniform vec2 WindowSize;

varying vec2 vTexCoord;
varying vec4 vColor;
varying float vType;

void main() {
    vec2 pos = Vertex / WindowSize * vec2(2.0, -2.0) + vec2(-1.0, 1.0);
    gl_Position = vec4(pos, 0, 1);
    vTexCoord = TexCoord;
    vColor = Color;
    vType = Type;
}
` + "\x00"

var fragmentSourceCode = `
varying vec2 vTexCoord;
varying vec4 vColor;
varying float vType;

void main() {
    if (vType == 0.0) {
        gl_FragColor = vColor;
    } else {
        gl_FragColor = vec4(1, 1, 0, 1);
    }
}
` + "\x00"
