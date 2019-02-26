package frontend

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/go-gl/gl/v4.1-core/gl"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

const (
	ppem    = 74
	width   = 96
	height  = 144
	originX = 0
	originY = 144
)

type Character struct {
	Texture uint32
	Image   *image.Alpha
}

type Font struct {
	characters map[rune]*Character
	font       *sfnt.Font
	Width      int
	Height     int
}

func NewFont(width, height int) (*Font, error) {
	f := &Font{
		characters: make(map[rune]*Character),
		Width:      width,
		Height:     height,
	}

	font, err := sfnt.Parse(gomono.TTF)
	if err != nil {
		return f, fmt.Errorf("could not parse font file: %v", err)
	}

	f.font = font

	return f, nil
}

func (f *Font) GetCharacter(r rune) (*Character, error) {
	var character *Character
	var err error

	character, ok := f.characters[r]
	if !ok {
		character, err = f.createCharacter(r)
		if err != nil {
			return character, fmt.Errorf("could not create character: %v", err)
		}
		f.characters[r] = character
	}
	return character, nil
}

func (f *Font) createCharacter(char rune) (*Character, error) {
	var character *Character
	var b sfnt.Buffer

	x, err := f.font.GlyphIndex(&b, char)
	if err != nil {
		return character, fmt.Errorf("could not get glyph index: %v", err)
	}
	if x == 0 {
		return character, fmt.Errorf("no glyph index found for rune '%s'", char)
	}

	segments, err := f.font.LoadGlyph(&b, x, fixed.I(ppem), nil)
	if err != nil {
		return character, fmt.Errorf("could not load glyph: %v", err)
	}

	r := vector.NewRasterizer(width, height)
	r.DrawOp = draw.Src
	for _, seg := range segments {
		// The divisions by 64 below is because the seg.Args values have type
		// fixed.Int26_6, a 26.6 fixed point number, and 1<<6 == 64.
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			r.MoveTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
			)
		case sfnt.SegmentOpLineTo:
			r.LineTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
			)
		case sfnt.SegmentOpQuadTo:
			r.QuadTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
				originX+float32(seg.Args[1].X)/64,
				originY+float32(seg.Args[1].Y)/64,
			)
		case sfnt.SegmentOpCubeTo:
			r.CubeTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
				originX+float32(seg.Args[1].X)/64,
				originY+float32(seg.Args[1].Y)/64,
				originX+float32(seg.Args[2].X)/64,
				originY+float32(seg.Args[2].Y)/64,
			)
		}
	}

	img := image.NewAlpha(image.Rect(0, 0, width, height))
	r.Draw(img, img.Bounds(), image.Opaque, image.Point{})

	var texture uint32
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(img.Bounds().Size().X), int32(img.Bounds().Size().Y), 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	// set texture options
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	character = &Character{
		Texture: texture,
		Image:   img,
	}

	return character, nil
}
