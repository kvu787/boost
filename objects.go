package main

import (
	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

type input_s struct {
	isMousePressed bool
	mousePosition  Vector
}

type transform_s struct {
	position     Vector
	velocity     Vector
	acceleration Vector
}

type circle_s struct {
	radius           float32
	outlineThickness float32
	rotation         float32
	fillColor        sf.Color
	outlineColor     sf.Color
}

func (c circle_s) GetDrawer() *sf.CircleShape {
	circle, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	circle.SetRadius(c.radius)
	circle.SetOrigin(NewCartesian(float64(c.radius/2), float64(c.radius/2)).ToSFMLVector2f())
	circle.SetOutlineThickness(c.outlineThickness)
	circle.SetRotation(c.rotation)
	circle.SetFillColor(c.fillColor)
	circle.SetOutlineColor(c.outlineColor)
	return circle
}

type player_s struct {
	transform transform_s
	circle    circle_s
}

type asteroid_s struct {
	transform transform_s
	circle_s  circle_s
}

func (t transform_s) applyAcceleration(duration float64) transform_s {
	newTransform := transform_s{}
	newTransform.acceleration = NewZeroVector()
	newTransform.velocity = t.velocity.Add(t.acceleration.Mul(duration))
	newTransform.position = t.position.Add(newTransform.velocity.Mul(duration))
	return newTransform
}
