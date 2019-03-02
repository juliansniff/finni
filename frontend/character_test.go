package frontend

import "testing"

func TestCharacter_SlideX(t *testing.T) {
	c := &Character{
		Vertices: [24]float32{},
	}

	c.SlideX(1)

	expected := [24]float32{
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
		1, 0, 0, 0,
	}

	for i := range expected {
		if c.Vertices[i] != expected[i] {
			t.Errorf("expected: %v, got: %v", expected, c.Vertices)
		}
	}
}

func TestCharacter_SlideY(t *testing.T) {
	c := &Character{
		Vertices: [24]float32{},
	}

	c.SlideY(1)

	expected := [24]float32{
		0, -1, 0, 0,
		0, -1, 0, 0,
		0, -1, 0, 0,
		0, -1, 0, 0,
		0, -1, 0, 0,
		0, -1, 0, 0,
	}

	for i := range expected {
		if c.Vertices[i] != expected[i] {
			t.Errorf("expected: %v, got: %v", expected, c.Vertices)
		}
	}
}
