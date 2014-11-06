package objects

import (
	. "bitbucket.org/kvu787/boost/lib/vector"
	"time"
)

type Transform_s struct {
	Position     Vector
	Velocity     Vector
	Acceleration Vector
}

func NewZeroTransform() Transform_s {
	return Transform_s{
		NewZeroVector(),
		NewZeroVector(),
		NewZeroVector(),
	}
}

func (t Transform_s) Act(duration time.Duration) Transform_s {
	newTransform := t
	velocityDelta := newTransform.Acceleration.Mul(duration.Seconds())
	newTransform.Velocity = newTransform.Velocity.Add(velocityDelta)
	positionDelta := newTransform.Velocity.Mul(duration.Seconds())
	newTransform.Position = newTransform.Position.Add(positionDelta)
	return newTransform
}
