package objects

type Asteroid_s struct {
	Transform        Transform_s
	Circle           Circle_s
	RenderProperties RenderProperties_s
}

func (a Asteroid_s) GetCircleShape() CircleShape_s {
	return CircleShape_s{a.Transform, a.Circle}
}
