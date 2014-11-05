package angle

import (
	"fmt"
	"math"
)

type Angle interface {
	GetRadians() float64
	GetDegrees() float64
	Add(angle Angle) Angle
	Sub(angle Angle) Angle
	Mul(constant float64) Angle
}

type angleType float64

func NewZeroAngle() Angle {
	return angleType(0)
}

func NewRadians(radians float64) Angle {
	return angleType(radians)
}

func NewDegrees(degrees float64) Angle {
	return NewRadians(degrees * ((2 * math.Pi) / 360.0))
}

func (a angleType) GetRadians() float64 {
	return float64(a)
}

func (a angleType) GetDegrees() float64 {
	return a.GetRadians() * (360.0 / (2 * math.Pi))
}

func (a angleType) Add(other Angle) Angle {
	return NewRadians(a.GetRadians() + other.GetRadians())
}

func (a angleType) Sub(other Angle) Angle {
	return NewRadians(a.GetRadians() - other.GetRadians())
}

func (a angleType) Mul(constant float64) Angle {
	return NewRadians(a.GetRadians() * constant)
}

func (a angleType) String() string {
	return fmt.Sprintf("Angle[ radians: %.2f, degrees: %.1f ]", a.GetRadians(), a.GetDegrees())
}
