package objects

import (
	"time"
)

type BumpTimeout_s struct {
	IsInBumpTimeout bool
	LastBumpTime    time.Time
}
