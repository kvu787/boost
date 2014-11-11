package objects

import (
	v "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

type RenderProperties_s struct {
	OutlineThickness float64
	Rotation         float64
	FillColor        sf.Color
	OutlineColor     sf.Color
	Scale            v.Vector
}
