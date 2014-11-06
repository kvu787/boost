package objects

type Player_s struct {
	Transform_s
	Circle_s
	RenderProperties_s
}

func (p Player_s) GetCircleShape() CircleShape_s {
	return CircleShape_s{p.Transform_s, p.Circle_s}
}
