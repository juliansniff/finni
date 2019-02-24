package frontend

import (
	"fmt"
	"image"
	"image/draw"

	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

const (
	ppem    = 32
	originX = 0
	originY = 36
)

type FontStore struct {
	bitmaps map[rune][]uint8
	font    *sfnt.Font
	Width   int
	Height  int
}

func NewFontStore(width, height int) (*FontStore, error) {
	fs := &FontStore{
		bitmaps: make(map[rune][]uint8),
		Width:   width,
		Height:  height,
	}

	font, err := sfnt.Parse(gomono.TTF)
	if err != nil {
		return fs, fmt.Errorf("could not parse font file: %v", err)
	}

	fs.font = font

	return fs, nil
}

func (fs *FontStore) GetBitmap(r rune) ([]uint8, error) {
	var b []uint8
	var err error

	b, ok := fs.bitmaps[r]
	if !ok {
		b, err = fs.createBitmap(r)
		if err != nil {
			return b, fmt.Errorf("could not create bitmap: %v", err)
		}
		fs.bitmaps[r] = b
	}
	return b, nil
}

func (fs *FontStore) createBitmap(char rune) ([]uint8, error) {
	var bitmap []uint8
	var b sfnt.Buffer

	x, err := fs.font.GlyphIndex(&b, char)
	if err != nil {
		return bitmap, fmt.Errorf("could not get glyph index: %v", err)
	}
	if x == 0 {
		return bitmap, fmt.Errorf("no glyph index found for rune '%s'", char)
	}

	segments, err := fs.font.LoadGlyph(&b, x, fixed.I(ppem), nil)
	if err != nil {
		return bitmap, fmt.Errorf("could not load glyph: %v", err)
	}

	r := vector.NewRasterizer(fs.Width, fs.Height)
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

	i := image.NewAlpha(image.Rect(0, 0, fs.Width, fs.Height))
	r.Draw(i, i.Bounds(), image.Opaque, image.Point{})

	return i.Pix, nil
}
