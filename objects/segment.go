package objects

import (
	v "bitbucket.org/kvu787/boost/lib/vector"
)

type Segment_s struct {
	End1, End2 v.Vector
}

func (s Segment_s) GetLength() float64 {
	return s.End1.Sub(s.End2).GetMagnitude()
}

func (s Segment_s) GetMidpoint() v.Vector {
	return s.End1.Add(s.End2.Sub(s.End1).Mul(0.5))
}
