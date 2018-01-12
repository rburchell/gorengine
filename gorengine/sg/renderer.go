package sg

import (
    "fmt"
    "github.com/go-gl/gl/v2.1/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
    "strings"
)

type Renderer struct {

    program uint32;
    vertexBuffer uint32;
    indexBuffer uint32;

    window *glfw.Window
}


func CreateRenderer() *Renderer {
    fmt.Println("Creating renderer")
    if err := glfw.Init(); err != nil {
        fmt.Println("Failed to initialize GLFW...")
        return nil
    }

    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 2)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)

    var renderer Renderer = Renderer{}
    var err error;

    renderer.window, err = glfw.CreateWindow(800, 600, "gorengine", nil, nil)
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


func renderPrepass(node Node) uint32 {
    var count uint32 = 0;
    if _, ok := node.(RectangleNode); ok {
        count = 1
    }

    for _, child := range node.GetChildren() {
        count += renderPrepass(child)
    }
    return count
}

func renderBuildBuffers(node Node, vertices []float32, vertexOffset uint32, indices []uint16, indexOffset uint32) {
    if rn, ok := node.(RectangleNode); ok {
        var firstIndex uint16 = uint16(vertexOffset / 9);
        indices[indexOffset + 0] = firstIndex + 0
        indices[indexOffset + 1] = firstIndex + 1
        indices[indexOffset + 2] = firstIndex + 2
        indices[indexOffset + 3] = firstIndex + 1
        indices[indexOffset + 4] = firstIndex + 3
        indices[indexOffset + 5] = firstIndex + 2
        indexOffset += 6;
        vertices[vertexOffset +  0 + 0] = rn.X
        vertices[vertexOffset +  0 + 1] = rn.Y
        vertices[vertexOffset +  0 + 2] = 0
        vertices[vertexOffset +  0 + 3] = 0
        vertices[vertexOffset +  0 + 4] = rn.R
        vertices[vertexOffset +  0 + 5] = rn.G
        vertices[vertexOffset +  0 + 6] = rn.B
        vertices[vertexOffset +  0 + 7] = rn.A
        vertices[vertexOffset +  0 + 8] = 0.0
        vertices[vertexOffset +  9 + 0] = rn.X + rn.W
        vertices[vertexOffset +  9 + 1] = rn.Y
        vertices[vertexOffset +  9 + 2] = 0
        vertices[vertexOffset +  9 + 3] = 0
        vertices[vertexOffset +  9 + 4] = rn.R
        vertices[vertexOffset +  9 + 5] = rn.G
        vertices[vertexOffset +  9 + 6] = rn.B
        vertices[vertexOffset +  9 + 7] = rn.A
        vertices[vertexOffset +  9 + 8] = 0.0
        vertices[vertexOffset + 18 + 0] = rn.X
        vertices[vertexOffset + 18 + 1] = rn.Y + rn.H
        vertices[vertexOffset + 18 + 2] = 0
        vertices[vertexOffset + 18 + 3] = 0
        vertices[vertexOffset + 18 + 4] = rn.R
        vertices[vertexOffset + 18 + 5] = rn.G
        vertices[vertexOffset + 18 + 6] = rn.B
        vertices[vertexOffset + 18 + 7] = rn.A
        vertices[vertexOffset + 18 + 8] = 0.0
        vertices[vertexOffset + 27 + 0] = rn.X + rn.W
        vertices[vertexOffset + 27 + 1] = rn.Y + rn.H
        vertices[vertexOffset + 27 + 2] = 0
        vertices[vertexOffset + 27 + 3] = 0
        vertices[vertexOffset + 27 + 4] = rn.R
        vertices[vertexOffset + 27 + 5] = rn.G
        vertices[vertexOffset + 27 + 6] = rn.B
        vertices[vertexOffset + 27 + 7] = rn.A
        vertices[vertexOffset + 27 + 8] = 0.0
        vertexOffset += 36;
    }
    for _, child := range node.GetChildren() {
        renderBuildBuffers(child, vertices, vertexOffset, indices, indexOffset)
    }
}


func (this *Renderer) Render(root Node) {
    gl.Clear(gl.COLOR_BUFFER_BIT)

    count := renderPrepass(root)
    vertexCount := count * 4
    indexCount := count * 6
    vertexBufferSize := vertexCount * 9
    indexBufferSize := indexCount

    vertexData := make([]float32, vertexBufferSize)
    indexData := make([]uint16, indexBufferSize)

    renderBuildBuffers(root, vertexData, 0, indexData, 0)

    for i:=0; i<int(vertexCount); i++ {
        fmt.Print("vertex: ", i, ": ");
        for v:=0; v<9; v++ {
            fmt.Print(vertexData[i * 9 + v], ", ")
        }
        fmt.Println()
    }
    for i:=0; i<int(count); i++ {
        fmt.Printf("index %d: %d %d %d %d %d %d\n", i,
                    indexData[i * 6 + 0],
                    indexData[i * 6 + 1],
                    indexData[i * 6 + 2],
                    indexData[i * 6 + 3],
                    indexData[i * 6 + 4],
                    indexData[i * 6 + 5])
    }

    gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexBuffer);
    gl.BufferData(gl.ARRAY_BUFFER, int(vertexBufferSize * 4), gl.Ptr(vertexData), gl.STREAM_DRAW)
    // gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.indexBuffer);
    // gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indexBufferSize * 2), gl.Ptr(indexData), gl.STREAM_DRAW);

    gl.EnableVertexAttribArray(0);
    gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 9 * 4 * 4, gl.PtrOffset(0));

    gl.UseProgram(this.program)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(vertexCount));

    gl.UseProgram(0)
    gl.DisableVertexAttribArray(0);

    this.window.SwapBuffers()
    glfw.PollEvents()
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

varying vec2 vTexCoord;
varying vec4 vColor;
varying float vType;

void main() {
    gl_Position = vec4(Vertex, 0, 1);
    vTexCoord = TexCoord;
    vColor = Color;
    vType = Type;
}
`

var fragmentSourceCode =
`
varying vec2 vTexCoord;
varying vec4 vColor;
varying float vType;

void main() {
    if (vType == 0.0) {
        gl_FragColor = vec4(1, 0, 0, 1);
    } else {
        gl_FragColor = vec4(1, 0, 1, 1);
    }
}
`


