package objects

type CircleRenderer_s struct {
	Circle_s
	RenderProperties_s
}

type CircleShape_s struct {
	Transform_s
	Circle_s
}

type Circle_s struct {
	Radius float64
}
