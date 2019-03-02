package frontend

type Character struct {
	Vertices [24]float32
	Texture  uint32
}

func (c *Character) SlideX(amount float32) {
	for i := 0; i < 6; i++ {
		c.Vertices[i*4] += amount
	}
}

func (c *Character) SlideY(amount float32) {
	for i := 0; i < 6; i++ {
		c.Vertices[i*4+1] -= amount
	}
}

func (c *Character) Scale(factor float32) {
	for i := 0; i < 6; i++ {
		for j := 0; j < 2; j++ {
			c.Vertices[i*4+j] *= factor
		}
	}
}
