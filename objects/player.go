package objects

type Player_s struct {
	Transform        Transform_s
	Circle           Circle_s
	RenderProperties RenderProperties_s
}

func (p Player_s) GetCircleShape() CircleShape_s {
	return CircleShape_s{p.Transform, p.Circle}
}
