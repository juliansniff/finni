package frontend

import (
	"fmt"
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

const (
	ppem    = 32
	width   = 24
	height  = 36
	originX = 0
	originY = 32
)

type Font struct {
	characters map[rune]Character
	font       *sfnt.Font
	metrics    font.Metrics
}

func NewFont() (*Font, error) {
	f := &Font{
		characters: make(map[rune]Character),
	}

	parsedFont, err := sfnt.Parse(gomono.TTF)
	if err != nil {
		return f, fmt.Errorf("could not parse font file: %v", err)
	}
	f.font = parsedFont

	metrics, err := parsedFont.Metrics(nil, ppem, font.HintingNone)
	if err != nil {
		return f, fmt.Errorf("could not get font metrics: %v", err)
	}
	f.metrics = metrics

	return f, nil
}

func (f *Font) GetCharacter(r rune) (Character, error) {
	var c Character
	var ok bool
	var err error

	c, ok = f.characters[r]
	if !ok {
		c, err = f.createCharacter(r)
		if err != nil {
			return c, fmt.Errorf("could not create character %v: %v", r, err)
		}
		f.characters[r] = c
	}

	return c, nil
}

func (f *Font) createCharacter(r rune) (Character, error) {
	var character Character

	x, err := f.font.GlyphIndex(nil, r)
	if err != nil {
		return character, fmt.Errorf("could not get glyph index: %v", err)
	}
	if x == 0 {
		return character, fmt.Errorf("no glyph index found for rune '%v'", r)
	}

	segments, err := f.font.LoadGlyph(nil, x, fixed.I(ppem), nil)
	if err != nil {
		return character, fmt.Errorf("could not load glyph: %v", err)
	}

	rasterizer := vector.NewRasterizer(width, height)
	rasterizer.DrawOp = draw.Src
	for _, seg := range segments {
		// The divisions by 64 below is because the seg.Args values have type
		// fixed.Int26_6, a 26.6 fixed point number, and 1<<6 == 64.
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			rasterizer.MoveTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
			)
		case sfnt.SegmentOpLineTo:
			rasterizer.LineTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
			)
		case sfnt.SegmentOpQuadTo:
			rasterizer.QuadTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
				originX+float32(seg.Args[1].X)/64,
				originY+float32(seg.Args[1].Y)/64,
			)
		case sfnt.SegmentOpCubeTo:
			rasterizer.CubeTo(
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
	rasterizer.Draw(img, img.Bounds(), image.Opaque, image.Point{})

	// create texture
	// var texture uint32
	// gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	// gl.GenTextures(1, &texture)
	// gl.BindTexture(gl.TEXTURE_2D, texture)
	// gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(img.Bounds().Size().X), int32(img.Bounds().Size().Y), 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// gl.BindTexture(gl.TEXTURE_2d, 0)
	// character.Texture = texture

	advance, err := f.font.GlyphAdvance(nil, x, ppem, font.HintingNone)
	if err != nil {
		return character, fmt.Errorf("could not get glyph advance: %v", err)
	}

	var xpos, ypos, h, w float32
	xpos = 0
	ypos = 0 - float32(f.metrics.Descent)/64
	h = float32(f.metrics.Height) / 64
	w = float32(advance) / 64
	character.Vertices = [12]float32{
		xpos, ypos + h,
		xpos, ypos,
		xpos + w, ypos,
		xpos, ypos + h,
		xpos + w, ypos,
		xpos + w, ypos + h,
	}

	return character, nil
}
