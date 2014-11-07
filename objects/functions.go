package objects

import (
	sf "bitbucket.org/kvu787/gosfml2"
)

func GetCircleShape(c Circle_s, r RenderProperties_s) *sf.CircleShape {
	circleShape, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	var radius float32 = float32(c.Radius)
	circleShape.SetRadius(radius - float32(r.OutlineThickness))
	circleShape.SetOrigin(sf.Vector2f{radius, radius})
	circleShape.SetOutlineThickness(float32(r.OutlineThickness))
	circleShape.SetRotation(float32(r.Rotation))
	circleShape.SetFillColor(r.FillColor)
	circleShape.SetOutlineColor(r.OutlineColor)
	return circleShape
}

func AreCirclesIntersecting(c1 CircleShape_s, c2 CircleShape_s, offset float64) bool {
	distance := c1.Transform_s.Position.Sub(c2.Transform_s.Position).GetMagnitude()
	radiusSum := c1.Circle_s.Radius + c2.Circle_s.Radius
	return distance+offset < radiusSum
}

func AreCircleSegmentIntersecting(s Segment_s, c CircleShape_s) bool {
	a := s.End2.Sub(s.End1)
	b := c.Position.Sub(s.End1)
	p := b.Projection(a)
	closestPointOnSegment := s.End1.Add(p)
	onLine := (closestPointOnSegment.GetX() < s.End1.GetX()) != (closestPointOnSegment.GetX() < s.End2.GetX())
	closeEnough := c.Radius > c.Position.Sub(closestPointOnSegment).GetMagnitude()
	return onLine && closeEnough
}

func GetCircleOverlap(c1 CircleShape_s, c2 CircleShape_s) float64 {
	sumRadius := c1.Radius + c2.Radius
	distance := c1.Position.Sub(c2.Position).GetMagnitude()
	return sumRadius - distance
}
