package objects

type Circle_s struct {
	Radius float64
}

// func (c circle_s) GetDrawer() *sf.CircleShape {
// 	circle, err := sf.NewCircleShape()
// 	if err != nil {
// 		panic(err)
// 	}
// 	circle.SetRadius(c.radius)
// 	circle.SetOrigin(NewCartesian(float64(c.radius/2), float64(c.radius/2)).ToSFMLVector2f())
// 	circle.SetOutlineThickness(c.outlineThickness)
// 	circle.SetRotation(c.rotation)
// 	circle.SetFillColor(c.fillColor)
// 	circle.SetOutlineColor(c.outlineColor)
// 	return circle
// }
