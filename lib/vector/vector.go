package vector

import (
	"fmt"
	"math"

	. "bitbucket.org/kvu787/boost/lib/angle"

	sf "bitbucket.org/kvu787/gosfml2"
)

type Vector interface {
	Dot(other Vector) float64
	Projection(other Vector) Vector
	GetX() float64
	GetY() float64
	Add(other Vector) Vector
	Sub(other Vector) Vector
	Mul(n float64) Vector
	Div(n float64) Vector
	GetMagnitude() float64
	SetMagnitude(magnitude float64)
	GetAngle() Angle
	SetAngle(angle Angle)
	ToSFMLVector2f() sf.Vector2f
}

type vectorStruct struct {
	x, y float64
}

func NewZeroVector() Vector {
	return NewCartesian(0, 0)
}

func NewCartesian(x, y float64) Vector {
	return &vectorStruct{x, y}
}

func NewPolar(magnitude float64, angle Angle) Vector {
	v := NewCartesian(0, 0)
	v.SetMagnitude(magnitude)
	v.SetAngle(angle)
	return v
}

func (v vectorStruct) Projection(other Vector) Vector {
	return v.Mul(v.Dot(other) / math.Pow(v.GetMagnitude(), 2))
}

func (v vectorStruct) Dot(other Vector) float64 {
	return v.GetX()*other.GetX() + v.GetY()*other.GetY()
}
func (v vectorStruct) GetX() float64 { return v.x }
func (v vectorStruct) GetY() float64 { return v.y }

func (v vectorStruct) Add(other Vector) Vector {
	return NewCartesian(v.GetX()+other.GetX(), v.GetY()+other.GetY())
}
func (v vectorStruct) Sub(other Vector) Vector { return v.Add(other.Mul(-1)) }
func (v vectorStruct) Mul(n float64) Vector    { return NewCartesian(v.GetX()*n, v.GetY()*n) }
func (v vectorStruct) Div(n float64) Vector    { return v.Mul(1 / n) }

func (v vectorStruct) GetMagnitude() float64 { return math.Sqrt(v.GetX()*v.GetX() + v.GetY()*v.GetY()) }

func (v *vectorStruct) SetMagnitude(magnitude float64) {
	angle := v.GetAngle()
	v.x = math.Cos(angle.GetRadians()) * magnitude
	v.y = math.Sin(angle.GetRadians()) * magnitude
}

func (v vectorStruct) GetAngle() Angle {
	return NewRadians(math.Atan2(v.GetY(), v.GetX()))
}

func (v *vectorStruct) SetAngle(angle Angle) {
	magnitude := v.GetMagnitude()
	v.x = math.Cos(angle.GetRadians()) * magnitude
	v.y = math.Sin(angle.GetRadians()) * magnitude
}

func (v vectorStruct) Distance(other Vector) float64 {
	return v.Sub(other).GetMagnitude()
}

func (v *vectorStruct) Direction(other Vector) Angle {
	return (other.Sub(v)).GetAngle()
}

func (v vectorStruct) ToSFMLVector2f() sf.Vector2f {
	return sf.Vector2f{float32(v.GetX()), float32(v.GetY())}
}

func (v vectorStruct) String() string {
	return fmt.Sprintf("vectorStruct[ x: %.1f, y: %.1f ]", v.GetX(), v.GetY())
	// return fmt.Sprintf("vectorStruct[ mag: %.1f, angle: %.1f ]", v.GetMagnitude(), v.GetAngle().GetDegrees())
}
