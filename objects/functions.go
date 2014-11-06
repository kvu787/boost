package objects

import (
	sf "bitbucket.org/kvu787/gosfml2"
)

func GetCircleShape(c Circle_s, r RenderProperties_s) *sf.CircleShape {
	circleShape, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	circleShape.SetRadius(float32(c.Radius))
	circleShape.SetOrigin(sf.Vector2f{float32(c.Radius) / 2, float32(c.Radius) / 2})
	circleShape.SetOutlineThickness(float32(r.OutlineThickness))
	circleShape.SetRotation(float32(r.Rotation))
	circleShape.SetFillColor(r.FillColor)
	circleShape.SetOutlineColor(r.OutlineColor)
	return circleShape
}
