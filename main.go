package main

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	p "bitbucket.org/kvu787/boost/lib/palette"
	v "bitbucket.org/kvu787/boost/lib/vector"
	o "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

func main() {
	// timing variables
	startTime := time.Now()
	durationPerFrame := time.Duration(int64(time.Second) / int64(FPS))

	// setup
	setupWindow()

	// setup game variables
	setupGameVariables()

	spawnInitialAsteroids(ASTEROID_INITIAL_SPAWN_COUNT)

	// main loop
	for WINDOW.IsOpen() {

		// vsync
		startTime = time.Now()

		input()  // update global INPUT
		update() // update global GAME_OBJECTS
		render() // write to WINDOW

		// sleep if frametime is short
		frametime := time.Since(startTime)
		framerate := 1.0 / frametime.Seconds()
		time.Sleep(time.Duration(int64(durationPerFrame) - int64(frametime)))
		fmt.Println("framerate: ", framerate)
	}
}

func setupWindow() {
	runtime.LockOSThread()
	WINDOW = sf.NewRenderWindow(
		sf.VideoMode{WINDOW_SIZE_X, WINDOW_SIZE_Y, 32},
		"boost",
		sf.StyleDefault,
		sf.DefaultContextSettings())
}

func setupGameVariables() {
	DURATION_PER_FRAME = time.Duration(int64(time.Second) / int64(FPS))
	CAMERA_OFFSET = v.NewCartesian(CAMERA_OFFSET_X, CAMERA_OFFSET_Y)
	INPUT = &o.Input_s{false, nil}
	PLAYER = &o.Player_s{
		o.Transform_s{
			v.NewZeroVector(),
			v.NewZeroVector(),
			v.NewZeroVector(),
		},
		o.Circle_s{5},
		o.RenderProperties_s{1, 0, p.BLUE, p.WHITE}}
	ASTEROIDS = list.New()
	ASTEROID_COLORS = []sf.Color{
		p.LIGHT_GRAY,
		p.GRAY,
		p.DARK_BROWN,
		p.LIGHT_BROWN,
	}
}

func spawnInitialAsteroids(n uint) {
	for n != 0 {
		// create asteroid
		color := ASTEROID_COLORS[rand.Intn(len(ASTEROID_COLORS))]
		velocityMagnitude := rand.Float64()*(ASTEROID_MAX_VELOCITY-ASTEROID_MIN_VELOCITY) + ASTEROID_MIN_VELOCITY
		radius := rand.Float64()*(ASTEROID_MAX_RADIUS-ASTEROID_MIN_RADIUS) + ASTEROID_MIN_RADIUS

		velocity := v.NewPolar(velocityMagnitude, rand.Float64()*2*math.Pi)
		newAsteroid := &o.Asteroid_s{
			o.Transform_s{
				v.NewZeroVector(),
				velocity,
				v.NewZeroVector(),
			},
			o.Circle_s{radius},
			o.RenderProperties_s{0, 0, color, p.WHITE}}

		// check if it intersects with anything
		for {
			// generate random position
			magnitude := float64(SPAWN_BOUNDARY) * rand.Float64()
			angle := 2 * math.Pi * rand.Float64()
			newAsteroid.Transform.Position = v.NewPolar(magnitude, angle)

			// check intersection with player
			if o.AreCirclesIntersecting(
				o.CircleShape_s{PLAYER.Transform, PLAYER.Circle},
				o.CircleShape_s{newAsteroid.Transform, newAsteroid.Circle}, -100) {
				continue
			}

			// check intersection with other asteroids
			isIntersectingWithOtherAsteroid := listAny(ASTEROIDS, func(i interface{}) bool {
				a := i.(*o.Asteroid_s)
				return o.AreCirclesIntersecting(
					o.CircleShape_s{a.Transform, a.Circle},
					o.CircleShape_s{newAsteroid.Transform, newAsteroid.Circle}, -5)
			})
			if isIntersectingWithOtherAsteroid {
				continue
			}

			break
		}
		ASTEROIDS.PushBack(newAsteroid)
		n--
	}
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
	INPUT.MousePosition = v.NewCartesian(float64(position.X), float64(position.Y))
}

func update() {
	fmt.Println("asteroids: ", ASTEROIDS.Len())
	// spawn asteroid
	if rand.Intn(int(ASTEROID_SPAWN_FREQUENCY)) == 0 {
		color := ASTEROID_COLORS[rand.Intn(len(ASTEROID_COLORS))]
		radian := v.DegreesToRadians(uint(rand.Intn(360)))
		velocityMagnitude := rand.Float64()*(ASTEROID_MAX_VELOCITY-ASTEROID_MIN_VELOCITY) + ASTEROID_MIN_VELOCITY
		velocityDegreeSpread := 20.0
		radius := rand.Float64()*(ASTEROID_MAX_RADIUS-ASTEROID_MIN_RADIUS) + ASTEROID_MIN_RADIUS

		position := v.NewPolar(SPAWN_BOUNDARY-10, radian)
		velocity := v.NewZeroVector().Sub(position)
		velocity.SetMagnitude(velocityMagnitude)
		velocity.SetAngle(
			velocity.GetAngle() -
				velocityDegreeSpread/2.0 +
				(velocityDegreeSpread * rand.Float64()))
		newAsteroid := &o.Asteroid_s{
			o.Transform_s{
				position,
				velocity,
				v.NewZeroVector(),
			},
			o.Circle_s{float64(radius)},
			o.RenderProperties_s{0, 0, color, p.WHITE}}
		ASTEROIDS.PushBack(newAsteroid)
	}

	// update player acceleration with user input
	if INPUT.IsMousePressed {
		camera := PLAYER.Transform.Position.Add(CAMERA_OFFSET)
		framedPlayerPosition := getFramedPosition(camera, PLAYER.Transform.Position)
		acceleration := framedPlayerPosition.Sub(INPUT.MousePosition)
		PLAYER.Transform.Acceleration = acceleration.Mul(SENSITIVITY)
	} else {
		PLAYER.Transform.Acceleration = v.NewZeroVector()
	}

	// update player transform
	PLAYER.Transform = PLAYER.Transform.Act(DURATION_PER_FRAME)

	// check if player is out of bounds
	if PLAYER.Transform.Position.GetMagnitude() >= float64(PLAYER_BOUNDARY) {
		PLAYER.Transform.Position = v.NewPolar(
			float64(PLAYER_BOUNDARY-PLAYER_RESET_DISTANCE),
			PLAYER.Transform.Position.GetAngle())

		PLAYER.Transform.Acceleration = v.NewZeroVector()

		PLAYER.Transform.Velocity = v.NewPolar(
			PLAYER_RESET_VELOCITY,
			PLAYER.Transform.Position.GetAngle()+math.Pi)
	}

	// remove asteroids that are out of bounds
	for e, next := ASTEROIDS.Front(), new(list.Element); e != nil; e = next {
		next = e.Next()
		asteroid := e.Value.(*o.Asteroid_s)
		if asteroid.Transform.Position.GetMagnitude() > SPAWN_BOUNDARY {
			ASTEROIDS.Remove(e)
		}
	}

	// update asteroid transforms
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*o.Asteroid_s)
		asteroid.Transform = asteroid.Transform.Act(DURATION_PER_FRAME)
	}

	// collide asteroids
	for e1 := ASTEROIDS.Front(); e1 != nil; e1 = e1.Next() {
		for e2 := e1.Next(); e2 != nil; e2 = e2.Next() {
			a1 := e1.Value.(*o.Asteroid_s)
			a2 := e2.Value.(*o.Asteroid_s)

			isIntersecting := o.AreCirclesIntersecting(
				a1.GetCircleShape(),
				a2.GetCircleShape(),
				3)

			if isIntersecting {
				// wimpy collision resolution
				a1.Transform.Velocity, a2.Transform.Velocity = resolveCollisionVelocities(a1.Transform, a2.Transform)

				// separate asteroids
				overlap := o.GetCircleOverlap(a1.GetCircleShape(), a2.GetCircleShape()) + 2
				displacement := a2.Transform.Position.Sub(a1.Transform.Position)
				displacement.SetMagnitude(overlap / 2)
				a2.Transform.Position = a2.Transform.Position.Add(displacement)
				displacement = displacement.Mul(-1)
				a1.Transform.Position = a1.Transform.Position.Add(displacement)
			}
		}

		// collide player
		a := e1.Value.(*o.Asteroid_s)
		isIntersecting := o.AreCirclesIntersecting(a.GetCircleShape(), PLAYER.GetCircleShape(), 1)
		if isIntersecting {
			PLAYER.Transform.Velocity, _ = resolveCollisionVelocities(PLAYER.Transform, a.Transform)
			PLAYER.Transform.Velocity.SetMagnitude(ASTEROID_BOUNCE_BACK_VELOCITY)
			temp := PLAYER.Transform.Position.Sub(a.Transform.Position)
			temp.SetMagnitude(a.Circle.Radius + PLAYER.Circle.Radius + 1)
			PLAYER.Transform.Position = a.Transform.Position.Add(temp)
		}
	}
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

func render() {
	WINDOW.Clear(p.BLACK)

	// get camera
	var camera v.Vector = PLAYER.Transform.Position.Add(CAMERA_OFFSET)

	// render player boundary
	boundaryCenterPosition := getFramedPosition(camera, v.NewZeroVector())
	c, err := sf.NewCircleShape()
	if err != nil {
		panic(err)
	}
	var pb float32 = float32(PLAYER_BOUNDARY + 10)
	c.SetPosition(boundaryCenterPosition.ToSFMLVector2f())
	c.SetRadius(pb)
	c.SetOrigin(sf.Vector2f{pb, pb})
	c.SetOutlineThickness(5)
	c.SetOutlineColor(p.WHITE)
	c.SetFillColor(p.TRANSPARENT)
	WINDOW.Draw(c, sf.DefaultRenderStates())

	// display player
	framedPlayerPosition := getFramedPosition(camera, PLAYER.Transform.Position)
	playerCircleShape := o.GetCircleShape(PLAYER.Circle, PLAYER.RenderProperties)
	playerCircleShape.SetPosition(framedPlayerPosition.ToSFMLVector2f())
	WINDOW.Draw(playerCircleShape, sf.DefaultRenderStates())

	// display asteroids
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*o.Asteroid_s)
		framedAsteroidPosition := getFramedPosition(camera, asteroid.Transform.Position)
		asteroidCircleShape := o.GetCircleShape(asteroid.Circle, asteroid.RenderProperties)
		asteroidCircleShape.SetPosition(framedAsteroidPosition.ToSFMLVector2f())
		WINDOW.Draw(asteroidCircleShape, sf.DefaultRenderStates())
	}

	WINDOW.Display()
}
