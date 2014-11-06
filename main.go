package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"
	. "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

func main() {
	// timing variables
	startTime := time.Now()
	durationPerFrame := time.Duration(int64(time.Second) / int64(FPS))

	// setup
	setup()

	// main loop
	for WINDOW.IsOpen() {

		// vsync
		startTime = time.Now()

		input()  // update global INPUT
		update() // update global GAME_OBJECTS
		render() // write to WINDOW

		// sleep if frametime is short
		time.Sleep(time.Duration(int64(durationPerFrame) - int64(time.Since(startTime))))
	}
}

func setup() {
	runtime.LockOSThread()
	WINDOW = sf.NewRenderWindow(
		sf.VideoMode{WINDOW_SIZE_X, WINDOW_SIZE_Y, 32},
		"boost",
		sf.StyleDefault,
		sf.DefaultContextSettings())
}

func input() {
	for event := WINDOW.PollEvent(); event != nil; event = WINDOW.PollEvent() {
		switch event.(type) {
		case sf.EventClosed:
			WINDOW.Close()
		case sf.EventMouseButtonPressed:
			INPUT.IsMousePressed = true
		case sf.EventMouseButtonReleased:
			INPUT.IsMousePressed = false
		}
	}
	position := sf.MouseGetPosition(WINDOW)
	INPUT.MousePosition = NewCartesian(float64(position.X), float64(position.Y))
}

func update() {
	fmt.Println(ASTEROIDS.Len())
	// spawn asteroid
	if rand.Intn(30) == 0 {
		color := ASTEROID_COLORS[rand.Intn(len(ASTEROID_COLORS))]
		radian := DegreesToRadians(uint(rand.Intn(360)))
		velocityMagnitude := rand.Intn(20) + 80
		velocityDegreeSpread := 30.0
		radius := rand.Intn(30) + 70

		position := NewPolar(SPAWN_BOUNDARY-10, float64(radian))
		velocity := NewZeroVector().Sub(position)
		velocity.SetMagnitude(float64(velocityMagnitude))
		velocity.SetAngle(
			velocity.GetAngle() -
				velocityDegreeSpread/2.0 +
				(velocityDegreeSpread * rand.Float64()))
		newAsteroid := &Asteroid_s{
			Transform_s{
				position,
				velocity,
				NewZeroVector(),
			},
			Circle_s{float64(radius)},
			RenderProperties_s{0, 0, color, palette.WHITE}}
		ASTEROIDS.PushBack(newAsteroid)
	}

	// update player acceleration with user input
	if INPUT.IsMousePressed {
		camera := PLAYER.Transform.Position.Add(CAMERA_OFFSET)
		framedPlayerPosition := getFramedPosition(camera, PLAYER.Transform.Position)
		acceleration := framedPlayerPosition.Sub(INPUT.MousePosition)
		PLAYER.Transform.Acceleration = acceleration
	} else {
		PLAYER.Transform.Acceleration = NewZeroVector()
	}

	// update player transform
	PLAYER.Transform = PLAYER.Transform.Act(DURATION_PER_FRAME)

	// check if player is out of bounds
	if PLAYER.Transform.Position.GetMagnitude() >= float64(PLAYER_BOUNDARY) {
		PLAYER.Transform.Position = NewPolar(
			float64(PLAYER_BOUNDARY-PLAYER_RESET_DISTANCE),
			PLAYER.Transform.Position.GetAngle())
		PLAYER.Transform.Acceleration = NewZeroVector()
		PLAYER.Transform.Velocity = NewZeroVector()
	}

	// remove asteroids that are out of bounds
	for e, next := ASTEROIDS.Front(), new(list.Element); e != nil; e = next {
		next = e.Next()
		asteroid := e.Value.(*Asteroid_s)
		if asteroid.Transform.Position.GetMagnitude() > SPAWN_BOUNDARY {
			ASTEROIDS.Remove(e)
		}
	}

	// update asteroid transforms
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*Asteroid_s)
		asteroid.Transform = asteroid.Transform.Act(DURATION_PER_FRAME)
	}

	// collide asteroids
	for e1 := ASTEROIDS.Front(); e1 != nil; e1 = e1.Next() {
		for e2 := e1.Next(); e2 != nil; e2 = e2.Next() {
			a1 := e1.Value.(*Asteroid_s)
			a2 := e2.Value.(*Asteroid_s)

			distance := a1.Transform.Position.Sub(a2.Transform.Position).GetMagnitude()
			sumRadius := a1.Circle.Radius + a2.Circle.Radius
			isIntersecting := distance-5 < sumRadius
			if isIntersecting {
				// wimpy collision resolution
				resolve := func(x1, x2 Transform_s) Vector {
					fromx1tox2 := x2.Position.Sub(x1.Position)
					proj := x1.Velocity.Projection(fromx1tox2)
					rej := x1.Velocity.Rejection(fromx1tox2)
					proj = proj.Mul(-1)
					return proj.Add(rej)
				}
				a1.Transform.Velocity = resolve(a1.Transform, a2.Transform)
				a2.Transform.Velocity = resolve(a2.Transform, a1.Transform)
			}
		}
	}
}

func render() {
	WINDOW.Clear(palette.BLACK)

	// get camera
	var camera Vector = PLAYER.Transform.Position.Add(CAMERA_OFFSET)

	// render player boundary
	boundaryCenterPosition := getFramedPosition(camera, NewZeroVector())
	c, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	var pb float32 = float32(PLAYER_BOUNDARY + 10)
	c.SetPosition(boundaryCenterPosition.ToSFMLVector2f())
	c.SetRadius(pb)
	c.SetOrigin(sf.Vector2f{pb, pb})
	c.SetOutlineThickness(5)
	c.SetOutlineColor(palette.WHITE)
	c.SetFillColor(palette.TRANSPARENT)
	WINDOW.Draw(c, sf.DefaultRenderStates())

	// display player
	framedPlayerPosition := getFramedPosition(camera, PLAYER.Transform.Position)
	playerCircleShape := GetCircleShape(PLAYER.Circle, PLAYER.RenderProperties)
	playerCircleShape.SetPosition(framedPlayerPosition.ToSFMLVector2f())
	WINDOW.Draw(playerCircleShape, sf.DefaultRenderStates())

	// display asteroids
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*Asteroid_s)
		framedAsteroidPosition := getFramedPosition(camera, asteroid.Transform.Position)
		asteroidCircleShape := GetCircleShape(asteroid.Circle, asteroid.RenderProperties)
		asteroidCircleShape.SetPosition(framedAsteroidPosition.ToSFMLVector2f())
		WINDOW.Draw(asteroidCircleShape, sf.DefaultRenderStates())
	}

	WINDOW.Display()
}
