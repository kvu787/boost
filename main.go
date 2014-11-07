package main

import (
	"container/list"
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
		time.Sleep(time.Duration(int64(durationPerFrame) - int64(frametime)))
		// fmt.Println("framerate: ", 1.0 / frametime.Seconds())
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
	SLIPS = list.New()

	// pushFrontAll(ASTEROIDS,
	// 	&o.Asteroid_s{
	// 		o.Transform_s{
	// 			v.NewCartesian(0, -100),
	// 			v.NewZeroVector(),
	// 			v.NewZeroVector()},
	// 		o.Circle_s{50},
	// 		o.RenderProperties_s{2, 0, p.WHITE, p.GRAY}},
	// 	&o.Asteroid_s{
	// 		o.Transform_s{
	// 			v.NewCartesian(100, 0),
	// 			v.NewZeroVector(),
	// 			v.NewZeroVector()},
	// 		o.Circle_s{50},
	// 		o.RenderProperties_s{2, 0, p.WHITE, p.GRAY}})

	ASTEROID_COLORS = []sf.Color{
		p.LIGHT_GRAY,
		p.GRAY,
		p.DARK_BROWN,
		p.LIGHT_BROWN}

	WINDOW_DIAGNOL_LENGTH = math.Sqrt(math.Pow(float64(WINDOW_SIZE_X), 2) + math.Pow(float64(WINDOW_SIZE_Y), 2))

	CURRENT_BOOST = MAX_BOOST
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
			newAsteroid.Position = v.NewPolar(magnitude, angle)

			// check intersection with player
			if o.AreCirclesIntersecting(
				o.CircleShape_s{PLAYER.Transform_s, PLAYER.Circle_s},
				o.CircleShape_s{newAsteroid.Transform_s, newAsteroid.Circle_s}, -100) {
				continue
			}

			// check intersection with other asteroids
			isIntersectingWithOtherAsteroid := listAny(ASTEROIDS, func(i interface{}) bool {
				a := i.(*o.Asteroid_s)
				return o.AreCirclesIntersecting(
					o.CircleShape_s{a.Transform_s, a.Circle_s},
					o.CircleShape_s{newAsteroid.Transform_s, newAsteroid.Circle_s}, -5)
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

	// spawn asteroid
	if SHOULD_SPAWN_ASTEROIDS {
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
	}

	// update player acceleration with user input
	if INPUT.IsMousePressed && CURRENT_BOOST > 0.0 {
		// apply accleration, decrease boost
		camera := PLAYER.Position.Add(CAMERA_OFFSET)
		framedPlayerPosition := getFramedPosition(camera, PLAYER.Position)
		acceleration := framedPlayerPosition.Sub(INPUT.MousePosition)
		PLAYER.Acceleration = acceleration.Mul(SENSITIVITY)
		CURRENT_BOOST -= BOOST_BURN * DURATION_PER_FRAME.Seconds()
	} else {
		PLAYER.Acceleration = v.NewZeroVector()
		if !INPUT.IsMousePressed {
			CURRENT_BOOST = math.Min(CURRENT_BOOST+BOOST_REGENERATION*DURATION_PER_FRAME.Seconds(), MAX_BOOST)
		}
	}

	// update player transform
	PLAYER.Transform_s = PLAYER.Transform_s.Act(DURATION_PER_FRAME)

	// check if player is out of bounds
	if PLAYER.Position.GetMagnitude() >= float64(PLAYER_BOUNDARY) {
		PLAYER.Position = v.NewPolar(
			float64(PLAYER_BOUNDARY-PLAYER_RESET_DISTANCE),
			PLAYER.Position.GetAngle())

		PLAYER.Acceleration = v.NewZeroVector()

		PLAYER.Velocity = v.NewPolar(
			PLAYER_RESET_VELOCITY,
			PLAYER.Position.GetAngle()+math.Pi)
	}

	// remove asteroids that are out of bounds
	for e, next := ASTEROIDS.Front(), new(list.Element); e != nil; e = next {
		next = e.Next()
		asteroid := e.Value.(*o.Asteroid_s)
		if asteroid.Position.GetMagnitude() > SPAWN_BOUNDARY {
			ASTEROIDS.Remove(e)
		}
	}

	// update asteroid transforms
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*o.Asteroid_s)
		asteroid.Transform_s = asteroid.Transform_s.Act(DURATION_PER_FRAME)
	}

	// loop through asteroids
	for e1 := ASTEROIDS.Front(); e1 != nil; e1 = e1.Next() {

		// compare every asteroid with every other
		for e2 := e1.Next(); e2 != nil; e2 = e2.Next() {
			a1 := e1.Value.(*o.Asteroid_s)
			a2 := e2.Value.(*o.Asteroid_s)

			isIntersecting := o.AreCirclesIntersecting(
				a1.GetCircleShape(),
				a2.GetCircleShape(),
				3)

			if isIntersecting {
				// collide asteroids

				// wimpy collision resolution
				a1.Velocity, a2.Velocity = resolveCollisionVelocities(a1.Transform_s, a2.Transform_s)

				// separate asteroids
				overlap := o.GetCircleOverlap(a1.GetCircleShape(), a2.GetCircleShape()) + 2
				displacement := a2.Position.Sub(a1.Position)
				displacement.SetMagnitude(overlap / 2)
				a2.Position = a2.Position.Add(displacement)
				displacement = displacement.Mul(-1)
				a1.Position = a1.Position.Add(displacement)
			} else {
				// check slip
				segment := o.Segment_s{a1.Position, a2.Position}
				if segment.GetLength()-a1.Radius-a2.Radius < 200 || LIGHT_SPAWN_DURATION == 1337 {
					if o.AreCircleSegmentIntersecting(segment, PLAYER.GetCircleShape()) {
						width := SLIP_WIDTH_SCALING / segment.GetLength()

						slip := o.Slip_s{
							segment,
							width,
							p.RandomColor(),
							time.Duration(int64(SLIP_DURATION*100)) * time.Millisecond,
							time.Now()}
						SLIPS.PushFront(slip)
					}
				}
			}
		}

		// collide player
		a := e1.Value.(*o.Asteroid_s)
		isIntersecting := o.AreCirclesIntersecting(a.GetCircleShape(), PLAYER.GetCircleShape(), 1)
		if isIntersecting {
			PLAYER.Velocity, _ = resolveCollisionVelocities(PLAYER.Transform_s, a.Transform_s)
			PLAYER.Velocity.SetMagnitude(ASTEROID_BOUNCE_BACK_VELOCITY)
			temp := PLAYER.Position.Sub(a.Position)
			temp.SetMagnitude(a.Radius + PLAYER.Radius + 1)
			PLAYER.Position = a.Position.Add(temp)
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
	var camera v.Vector = PLAYER.Position.Add(CAMERA_OFFSET)

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

	// render player
	framedPlayerPosition := getFramedPosition(camera, PLAYER.Position)
	playerCircleShape := o.GetCircleShape(PLAYER.Circle_s, PLAYER.RenderProperties_s)
	playerCircleShape.SetPosition(framedPlayerPosition.ToSFMLVector2f())
	WINDOW.Draw(playerCircleShape, sf.DefaultRenderStates())

	// render asteroids
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*o.Asteroid_s)
		framedAsteroidPosition := getFramedPosition(camera, asteroid.Position)
		asteroidCircleShape := o.GetCircleShape(asteroid.Circle_s, asteroid.RenderProperties_s)
		asteroidCircleShape.SetPosition(framedAsteroidPosition.ToSFMLVector2f())
		WINDOW.Draw(asteroidCircleShape, sf.DefaultRenderStates())
	}

	// render slips
	for e, next := SLIPS.Front(), new(list.Element); e != nil; e = next {
		next = e.Next()
		slip := e.Value.(o.Slip_s)

		timeElapsed := time.Since(slip.TimeSpawned)
		if timeElapsed.Seconds() > slip.Duration.Seconds() {
			SLIPS.Remove(e)
		} else {
			transparencyRatio := (slip.Duration.Seconds() - timeElapsed.Seconds()) / slip.Duration.Seconds()
			color := slip.Color
			color.A = uint8(transparencyRatio * float64(255))
			r, err := sf.NewRectangleShape()
			if err != nil {
				panic(err)
			}
			r.SetPosition(getFramedPosition(camera, slip.Segment_s.GetMidpoint()).ToSFMLVector2f())
			r.SetSize(sf.Vector2f{float32(SPAWN_BOUNDARY) * 2, float32(slip.Width)})
			r.SetRotation(float32(v.RadiansToDegrees(slip.End1.Sub(slip.End2).GetAngle())))
			r.SetFillColor(color)
			r.SetOutlineThickness(0)
			r.SetOrigin(sf.Vector2f{float32(SPAWN_BOUNDARY), float32(slip.Width) / 2})
			WINDOW.Draw(r, sf.DefaultRenderStates())
		}
	}

	// render boost
	r, err := sf.NewRectangleShape()
	if err != nil {
		panic(err)
	}
	color := p.RED
	color.A = 200
	boostRatio := CURRENT_BOOST / MAX_BOOST
	r.SetPosition(v.NewZeroVector().ToSFMLVector2f())
	r.SetSize(v.NewCartesian(boostRatio*float64(WINDOW_SIZE_X), 10).ToSFMLVector2f())
	r.SetRotation(0)
	r.SetFillColor(color)
	r.SetOutlineThickness(0)
	WINDOW.Draw(r, sf.DefaultRenderStates())

	WINDOW.Display()
}
