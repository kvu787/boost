package objects

type Asteroid_s struct {
	Transform_s
	Circle_s
	RenderProperties_s
}

func (a Asteroid_s) GetCircleShape() CircleShape_s {
	return CircleShape_s{a.Transform_s, a.Circle_s}
}
