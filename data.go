package main

import (
	"container/list"
	"time"

	. "bitbucket.org/kvu787/boost/lib/vector"
	. "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

type GraphicsSettings_s struct {
	WindowX uint
	WindowY uint
	FpsLock uint
}

type PlayerSettings_s struct {
	Color           string
	Radius          float64
	ControlRadius   float64
	MaxAcceleration float64
}

type SandboxSettings_s struct {
	InitialSpawnCount        uint
	MaxCount                 uint
	SpawnFrequency           uint
	MinRadius                float64
	MaxRadius                float64
	MinVelocity              float64
	MaxVelocity              float64
	PlayerBouncebackVelocity float64
}

var GRAPHICS_SETTINGS GraphicsSettings_s
var PLAYER_SETTINGS PlayerSettings_s
var SANDBOX_SETTINGS SandboxSettings_s

const (
	PLAYER_BOUNDARY                float64 = 1000
	SPAWN_BOUNDARY                 float64 = 1400
	PLAYER_BOUNDARY_RESET_DISTANCE float64 = 3
	PLAYER_BOUNDARY_RESET_VELOCITY float64 = 200

	MAX_BOOST            float64 = 100
	BOOST_REGENERATION   float64 = 100
	BOOST_BURN           float64 = 100
	LIGHT_SPAWN_DURATION uint    = 0x529

	SLIP_DURATION      float64 = 8
	SLIP_WIDTH_SCALING float64 = 4500 // higher means wider
	SLIP_MAX_WIDTH     float64 = 300

	HEALTH_MAX    float64 = 100
	HEALTH_REGEN  float64 = 2
	HEALTH_DECAY  float64 = 1 // points per second
	HEALTH_DAMAGE float64 = 10

	ENDLESS bool = false
)

// variables that should be immutable
var (
	DURATION_PER_FRAME        time.Duration
	ASTEROID_COLORS           []sf.Color
	WINDOW_DIAGNOL_LENGTH     float64
	PLAYER_ACCELERATION_CURVE func(float64) float64
)

var (
	WINDOW                   *sf.RenderWindow
	FRAME                    Vector
	CAMERA_SHIFT             Vector // from center
	INPUT                    *Input_s
	PLAYER                   *Player_s
	ASTEROIDS                *list.List
	SLIPS                    *list.List
	CURRENT_BOOST            float64
	LAST_ASTEROID_SPAWN_TIME time.Time
	HEALTH_CURRENT           float64
)
