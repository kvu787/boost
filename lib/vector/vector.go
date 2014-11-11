package vector

import (
	"fmt"
	"math"

	sf "bitbucket.org/kvu787/gosfml2"
)

type Vector interface {
	GetX() float64
	GetY() float64
	GetMagnitude() float64
	GetAngle() float64

	Add(other Vector) Vector
	Sub(other Vector) Vector
	Mul(n float64) Vector
	Div(n float64) Vector

	Dot(other Vector) float64
	Projection(other Vector) Vector // parallel to other
	Rejection(other Vector) Vector  // perpendicular to other

	SetMagnitude(magnitude float64)
	SetAngle(angle float64)

	// Distance(other Vector) float64
	// Direction(other Vector) float64

	ToSFMLVector2f() sf.Vector2f
}

type vectorStruct struct {
	x, y float64
}

func RadiansToDegrees(radians float64) float64 {
	return (radians / (2.0 * math.Pi)) * 360.0
}

func DegreesToRadians(degrees uint) float64 {
	return (float64(degrees) / 360.0) * 2 * math.Pi
}

func NewUnitVector() Vector {
	return NewCartesian(1, 1)
}

func NewZeroVector() Vector {
	return NewCartesian(0, 0)
}

func NewCartesian(x, y float64) Vector {
	return &vectorStruct{x, y}
}

func NewPolar(magnitude float64, angle float64) Vector {
	v := NewCartesian(0, 0)
	v.SetMagnitude(magnitude)
	v.SetAngle(angle)
	return v
}

func (v vectorStruct) GetX() float64         { return v.x }
func (v vectorStruct) GetY() float64         { return v.y }
func (v vectorStruct) GetMagnitude() float64 { return math.Sqrt(v.GetX()*v.GetX() + v.GetY()*v.GetY()) }
func (v vectorStruct) GetAngle() float64 {
	return math.Atan2(v.GetY(), v.GetX())
}

func (v vectorStruct) Add(other Vector) Vector {
	return NewCartesian(v.GetX()+other.GetX(), v.GetY()+other.GetY())
}
func (v vectorStruct) Sub(other Vector) Vector { return v.Add(other.Mul(-1)) }
func (v vectorStruct) Mul(n float64) Vector    { return NewCartesian(v.GetX()*n, v.GetY()*n) }
func (v vectorStruct) Div(n float64) Vector    { return v.Mul(1 / n) }

func (v vectorStruct) Projection(other Vector) Vector {
	return other.Mul(v.Dot(other) / math.Pow(other.GetMagnitude(), 2))
}
func (v vectorStruct) Rejection(other Vector) Vector {
	return v.Sub(v.Projection(other))
}
func (v vectorStruct) Dot(other Vector) float64 {
	return v.GetX()*other.GetX() + v.GetY()*other.GetY()
}

func (v *vectorStruct) SetMagnitude(magnitude float64) {
	angle := v.GetAngle()
	v.x = math.Cos(angle) * magnitude
	v.y = math.Sin(angle) * magnitude
}
func (v *vectorStruct) SetAngle(angle float64) {
	magnitude := v.GetMagnitude()
	v.x = math.Cos(angle) * magnitude
	v.y = math.Sin(angle) * magnitude
}

func (v vectorStruct) Distance(other Vector) float64 {
	return v.Sub(other).GetMagnitude()
}
func (v *vectorStruct) Direction(other Vector) float64 {
	return (other.Sub(v)).GetAngle()
}

func (v vectorStruct) ToSFMLVector2f() sf.Vector2f {
	return sf.Vector2f{float32(v.GetX()), float32(v.GetY())}
}

func (v vectorStruct) String() string {
	return fmt.Sprintf("vectorStruct[ x: %.1f, y: %.1f ]", v.GetX(), v.GetY())
	// return fmt.Sprintf("vectorStruct[ mag: %.1f, angle: %.1f ]", v.GetMagnitude(), v.GetAngle().GetDegrees())
}
