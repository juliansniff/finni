package frontend

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var vertices = []float32{
	0.5, 0.5,
	0.5, -0.5,
	-0.5, -0.5,
	-0.5, 0.5,
}

type Buffer struct {
	Characters []Character
	VAO        uint32
	VBO        uint32
}

// NewBuffer creates a Buffer.
func NewBuffer(r rune, f *Font) (*Buffer, error) {
	b := &Buffer{
		Characters: make([]Character, 0),
	}

	c, err := f.GetCharacter(r)
	if err != nil {
		return b, fmt.Errorf("could not get character '%v': %v", r, err)
	}
	b.Characters = append(b.Characters, c)

	gl.GenVertexArrays(1, &b.VAO)
	gl.BindVertexArray(b.VAO)

	gl.GenBuffers(1, &b.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.VBO)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return b, nil
}

func (b *Buffer) Draw() {
	gl.BindVertexArray(b.VAO)

	for _, c := range b.Characters {
		gl.BindBuffer(gl.ARRAY_BUFFER, b.VBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.Vertices)*4, gl.Ptr(&c.Vertices[0]), gl.STATIC_DRAW)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(c.Vertices)))
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}
