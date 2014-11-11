package main

import (
	"container/list"
	"fmt"
	"math"

	v "bitbucket.org/kvu787/boost/lib/vector"
	o "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

func listAny(l *list.List, f func(interface{}) bool) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			return true
		}
	}
	return false
}

func listNew(elements ...interface{}) *list.List {
	result := list.New()
	for _, e := range elements {
		fmt.Println(e)
		result.PushBack(e)
	}
	return result
}

func pushFrontAll(l *list.List, objects ...interface{}) {
	for _, object := range objects {
		l.PushFront(object)
	}
}

func worldToFramePosition(frame, x v.Vector) v.Vector {
	return x.Sub(frame)
}

func frameToWorldPosition(frame, x v.Vector) v.Vector {
	return frame.Add(x)
}

func polynomial(exp, xScale, yScale float64) func(float64) float64 {
	return func(x float64) float64 {
		k := (yScale) / (math.Pow(xScale, exp))
		return k * (math.Pow(x, exp))
	}
}

// origin is center of the rectangle
// rotation is clockwise
// dimensions are {width, height}
func CreateRectangle(size v.Vector, rp o.RenderProperties_s) *sf.RectangleShape {
	r, err := sf.NewRectangleShape()
	if err != nil {
		panic(err)
	}
	r.SetSize(size.ToSFMLVector2f())

	r.SetOrigin(sf.Vector2f{float32(size.GetX() / 2.0), float32(size.GetY() / 2.0)})
	r.SetRotation(float32(v.RadiansToDegrees(rp.Rotation)))
	r.SetOutlineThickness(float32(rp.OutlineThickness))
	r.SetOutlineColor(rp.OutlineColor)
	r.SetFillColor(rp.FillColor)
	r.SetScale(rp.Scale.ToSFMLVector2f())
	return r
}

// origin is center of circle
func CreateCircle(radius float64, rp o.RenderProperties_s) *sf.CircleShape {
	c, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	r32 := float32(radius)
	c.SetRadius(r32)

	c.SetOrigin(sf.Vector2f{r32, r32})
	c.SetRotation(float32(rp.Rotation))
	c.SetOutlineThickness(float32(rp.OutlineThickness))
	c.SetOutlineColor(rp.OutlineColor)
	c.SetFillColor(rp.FillColor)
	c.SetScale(rp.Scale.ToSFMLVector2f())

	return c
}

func resolveCollisionVelocities(t1, t2 o.Transform_s) (v.Vector, v.Vector) {
	resolve := func(x1, x2 o.Transform_s) v.Vector {
		fromx1tox2 := x2.Position.Sub(x1.Position)
		proj := x1.Velocity.Projection(fromx1tox2)
		rej := x1.Velocity.Rejection(fromx1tox2)
		proj = proj.Mul(-1)
		return proj.Add(rej)
	}
	return resolve(t1, t2), resolve(t2, t1)
}
