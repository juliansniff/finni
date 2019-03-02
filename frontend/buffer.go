package frontend

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Buffer struct {
	Characters []Character
	VAO        uint32
	VBO        uint32
}

// NewBuffer creates a Buffer.
func NewBuffer(path string, f *Font) (*Buffer, error) {
	b := &Buffer{
		Characters: make([]Character, 0),
	}

	fBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return b, fmt.Errorf("could not open file '%s': %v", path, err)
	}

	var xSlide, ySlide float32
	for i := range fBytes {
		r := rune(fBytes[i])

		c, err := f.GetCharacter(r)
		if err != nil {
			log.Printf("could not get character '%v': %v", r, err)
			continue
		}

		c.Scale(0.1)
		if r == '\n' || r == '\r' {
			ySlide += 0.06
			xSlide = 0
		} else {
			xSlide += 0.025
		}

		c.SlideX(xSlide)
		c.SlideY(ySlide)
		b.Characters = append(b.Characters, c)
	}

	gl.GenVertexArrays(1, &b.VAO)
	gl.BindVertexArray(b.VAO)

	gl.GenBuffers(1, &b.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.VBO)

	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return b, nil
}

func (b *Buffer) Draw() {
	gl.BindVertexArray(b.VAO)
	gl.ActiveTexture(gl.TEXTURE0)

	for _, c := range b.Characters {
		gl.BindTexture(gl.TEXTURE_2D, c.Texture)
		gl.BindBuffer(gl.ARRAY_BUFFER, b.VBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(c.Vertices)*4, gl.Ptr(&c.Vertices[0]), gl.STATIC_DRAW)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(c.Vertices)))
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
