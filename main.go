package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	p "bitbucket.org/kvu787/boost/lib/palette"
	v "bitbucket.org/kvu787/boost/lib/vector"
	o "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

func main() {
	// read config files
	func(graphicsPath, playerPath, levelPath string) {
		bs, err := ioutil.ReadFile(graphicsPath)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, &GRAPHICS_SETTINGS)

		bs, err = ioutil.ReadFile(playerPath)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, &PLAYER_SETTINGS)

		bs, err = ioutil.ReadFile(levelPath)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bs, &SANDBOX_SETTINGS)
	}(os.Args[1], os.Args[2], os.Args[3])

	// timing variables
	startTime := time.Now()
	durationPerFrame := time.Duration(int64(time.Second) / int64(GRAPHICS_SETTINGS.FpsLock))

	// setup
	setupWindow()

	// setup game variables
	setupGameVariables()

	spawnInitialAsteroids(SANDBOX_SETTINGS.InitialSpawnCount)

	// main loop
	fmt.Println("hello my name is major tom")
	for WINDOW.IsOpen() {

		// vsync
		startTime = time.Now()

		input()  // update global INPUT
		update() // update global GAME_OBJECTS
		render() // write to WINDOW

		// sleep if frametime is short
		frametime := time.Since(startTime)
		time.Sleep(time.Duration(int64(durationPerFrame) - int64(frametime)))
	}
}

func setupWindow() {
	runtime.LockOSThread()
	WINDOW = sf.NewRenderWindow(
		sf.VideoMode{GRAPHICS_SETTINGS.WindowX, GRAPHICS_SETTINGS.WindowY, 32},
		"boost",
		sf.StyleDefault,
		sf.DefaultContextSettings())
}

func setupGameVariables() {
	DURATION_PER_FRAME = time.Duration(int64(time.Second) / int64(GRAPHICS_SETTINGS.FpsLock))
	FRAME = v.NewCartesian((float64(GRAPHICS_SETTINGS.WindowX))/-2.0, (float64(GRAPHICS_SETTINGS.WindowY))/-2.0)
	INPUT = &o.Input_s{false, nil}
	CAMERA_SHIFT = v.NewZeroVector()

	PLAYER = &o.Player_s{
		o.Transform_s{
			v.NewZeroVector(),
			v.NewZeroVector(),
			v.NewZeroVector(),
		},
		o.Circle_s{6},
		o.RenderProperties_s{0, 0, p.BLUE, p.WHITE, v.NewUnitVector()}}

	ASTEROIDS = list.New()
	SLIPS = list.New()

	ASTEROID_COLORS = []sf.Color{p.WHITE}

	// ASTEROID_COLORS = []sf.Color{
	// 	p.LIGHT_GRAY,
	// 	p.GRAY,
	// 	p.DARK_BROWN,
	// 	p.LIGHT_BROWN}

	WINDOW_DIAGNOL_LENGTH = math.Sqrt(math.Pow(float64(GRAPHICS_SETTINGS.WindowX), 2) + math.Pow(float64(GRAPHICS_SETTINGS.WindowY), 2))

	CURRENT_BOOST = MAX_BOOST

	PLAYER_ACCELERATION_CURVE = polynomial(2, PLAYER_SETTINGS.ControlRadius, PLAYER_SETTINGS.MaxAcceleration)

	LAST_ASTEROID_SPAWN_TIME = time.Now()

	HEALTH_CURRENT = HEALTH_MAX

	BUMP_TIMEOUT_DURATION = time.Second

	PLAYER_BUMP_TIMEOUT = &o.BumpTimeout_s{false, time.Now()}
}

func spawnInitialAsteroids(n uint) {
	for n != 0 {
		// create asteroid
		color := ASTEROID_COLORS[rand.Intn(len(ASTEROID_COLORS))]
		velocityMagnitude := rand.Float64()*(SANDBOX_SETTINGS.MaxVelocity-SANDBOX_SETTINGS.MinVelocity) + SANDBOX_SETTINGS.MinVelocity
		radius := rand.Float64()*(SANDBOX_SETTINGS.MaxRadius-SANDBOX_SETTINGS.MinRadius) + SANDBOX_SETTINGS.MinRadius

		velocity := v.NewPolar(velocityMagnitude, rand.Float64()*2*math.Pi)
		newAsteroid := &o.Asteroid_s{
			o.Transform_s{
				v.NewZeroVector(),
				velocity,
				v.NewZeroVector(),
			},
			o.Circle_s{radius},
			o.RenderProperties_s{0, 0, color, p.WHITE, v.NewCartesian(1, 1)}}

		// check if it intersects with anything
		for {
			// generate random position
			magnitude := float64(SPAWN_BOUNDARY) * rand.Float64()
			angle := 2 * math.Pi * rand.Float64()
			newAsteroid.Position = v.NewPolar(magnitude, angle)

			// check intersection with player
			if o.AreCirclesIntersecting(
				PLAYER.GetCircleShape(),
				newAsteroid.GetCircleShape(), -100) {
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
	fmt.Println(ASTEROIDS.Len())
	// spawn asteroid
	func(spawnFrequency uint) {
		var spawnDelay float64 = 1.0 / float64(spawnFrequency)
		if (uint(ASTEROIDS.Len()) < SANDBOX_SETTINGS.MaxCount) && time.Since(LAST_ASTEROID_SPAWN_TIME).Seconds() >= spawnDelay {
			LAST_ASTEROID_SPAWN_TIME = time.Now()
			color := ASTEROID_COLORS[rand.Intn(len(ASTEROID_COLORS))]
			radian := v.DegreesToRadians(uint(rand.Intn(360)))
			velocityMagnitude := rand.Float64()*(SANDBOX_SETTINGS.MaxVelocity-SANDBOX_SETTINGS.MinVelocity) + SANDBOX_SETTINGS.MinVelocity
			velocityDegreeSpread := 20.0
			radius := rand.Float64()*(SANDBOX_SETTINGS.MaxRadius-SANDBOX_SETTINGS.MinRadius) + SANDBOX_SETTINGS.MinRadius

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
				o.RenderProperties_s{0, 0, color, p.WHITE, v.NewCartesian(1, 1)}}
			ASTEROIDS.PushBack(newAsteroid)
		}
	}(SANDBOX_SETTINGS.SpawnFrequency)

	// update player bump timeout
	func() {
		if PLAYER_BUMP_TIMEOUT.IsInBumpTimeout {
			if time.Since(PLAYER_BUMP_TIMEOUT.LastBumpTime).Seconds() > BUMP_TIMEOUT_DURATION.Seconds() {
				// turn off bump timeout
				PLAYER_BUMP_TIMEOUT.IsInBumpTimeout = false
			}
		}
	}()

	// update player acceleration
	func() {
		if INPUT.IsMousePressed && !PLAYER_BUMP_TIMEOUT.IsInBumpTimeout {
			worldMousePosition := frameToWorldPosition(FRAME, INPUT.MousePosition)
			displacement := PLAYER.Position.Sub(worldMousePosition)
			acceleration := v.NewPolar(
				PLAYER_ACCELERATION_CURVE(math.Min(displacement.GetMagnitude(), PLAYER_SETTINGS.ControlRadius)),
				displacement.GetAngle())
			PLAYER.Acceleration = acceleration
		} else {
			PLAYER.Acceleration = v.NewZeroVector()
		}
	}()

	// update camera shift
	func() {
		frameCenter := v.NewCartesian((float64(GRAPHICS_SETTINGS.WindowX))/2.0, (float64(GRAPHICS_SETTINGS.WindowY))/2.0)
		centerMouseDisplacement := INPUT.MousePosition.Sub(frameCenter)
		shiftFromMovement := centerMouseDisplacement.Mul(-0.1)
		goalCameraShift := shiftFromMovement

		// if INPUT.IsMousePressed {
		// 	shiftFromAcceleration := centerMouseDisplacement.Mul(0.3 * PLAYER.Acceleration.GetMagnitude() / PLAYER_MAX_ACCELERATION)
		// 	goalCameraShift = goalCameraShift.Add(shiftFromAcceleration)
		// }

		diff := goalCameraShift.Sub(CAMERA_SHIFT)
		if INPUT.IsMousePressed {
			CAMERA_SHIFT = CAMERA_SHIFT.Add(diff.Mul(0.2))
		} else {
			CAMERA_SHIFT = CAMERA_SHIFT.Add(diff.Mul(0.1))
		}
	}()

	// update player transform
	PLAYER.Transform_s = PLAYER.Transform_s.Act(DURATION_PER_FRAME)

	// check if player is out of bounds
	if PLAYER.Position.GetMagnitude() >= float64(PLAYER_BOUNDARY) {
		PLAYER.Position = v.NewPolar(
			float64(PLAYER_BOUNDARY-PLAYER_BOUNDARY_RESET_DISTANCE),
			PLAYER.Position.GetAngle())

		PLAYER.Acceleration = v.NewZeroVector()

		PLAYER.Velocity = v.NewPolar(
			PLAYER_BOUNDARY_RESET_VELOCITY,
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

	// update player health
	func() {
		if HEALTH_CURRENT < 0 && !ENDLESS {
			panic("gg")
		}
		HEALTH_CURRENT -= HEALTH_DECAY * DURATION_PER_FRAME.Seconds()
	}()

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
				// check slip if not bumped
				if !PLAYER_BUMP_TIMEOUT.IsInBumpTimeout {
					segment := o.Segment_s{a1.Position, a2.Position}
					if segment.GetLength()-a1.Radius-a2.Radius < SLIP_MAX_WIDTH || LIGHT_SPAWN_DURATION == 02471 {
						if o.AreCircleSegmentIntersecting(segment, PLAYER.GetCircleShape()) {
							width := SLIP_WIDTH_SCALING / segment.GetLength()

							slip := o.Slip_s{
								segment,
								width,
								p.RandomColor(),
								time.Duration(int64(SLIP_DURATION*100)) * time.Millisecond,
								time.Now()}
							SLIPS.PushFront(slip)

							// regen player health
							HEALTH_CURRENT = math.Min(HEALTH_CURRENT+HEALTH_REGEN, HEALTH_MAX)
						}
					}
				}
			}
		}

		// collide player with asteroids
		a := e1.Value.(*o.Asteroid_s)
		isIntersecting := o.AreCirclesIntersecting(a.GetCircleShape(), PLAYER.GetCircleShape(), 1)
		if isIntersecting {

			func() { // activate bump timeout
				// be invincible if in bump timeout
				if !PLAYER_BUMP_TIMEOUT.IsInBumpTimeout {
					HEALTH_CURRENT -= HEALTH_DAMAGE
					PLAYER_BUMP_TIMEOUT.IsInBumpTimeout = true
					PLAYER_BUMP_TIMEOUT.LastBumpTime = time.Now()
				}
			}()

			PLAYER.Velocity, _ = resolveCollisionVelocities(PLAYER.Transform_s, a.Transform_s)
			PLAYER.Velocity.SetMagnitude(SANDBOX_SETTINGS.PlayerBouncebackVelocity)
			temp := PLAYER.Position.Sub(a.Position)
			temp.SetMagnitude(a.Radius + PLAYER.Radius + 1)
			PLAYER.Position = a.Position.Add(temp)
		}
	}
}

func render() {
	WINDOW.Clear(p.BLACK)

	// update the frame with respect to player position
	FRAME = PLAYER.Position.Add(v.NewCartesian((float64(GRAPHICS_SETTINGS.WindowX))/-2.0, (float64(GRAPHICS_SETTINGS.WindowY))/-2.0)).Add(CAMERA_SHIFT)

	// render grid lines
	func(density int, width float64) {
		c := p.WHITE
		c.A = 70
		// vertical
		for i := 0; i <= int(SPAWN_BOUNDARY)*2; i += density {
			rp := o.RenderProperties_s{
				0,
				math.Pi * 0.5,
				p.SetAlpha(p.WHITE, 70),
				p.WHITE,
				v.NewUnitVector(),
			}
			r := CreateRectangle(v.NewCartesian(SPAWN_BOUNDARY*2, width), rp)
			r.SetPosition(
				worldToFramePosition(FRAME,
					v.NewCartesian(-SPAWN_BOUNDARY+float64(i), 0)).ToSFMLVector2f())
			WINDOW.Draw(r, sf.DefaultRenderStates())
		}
		// horizontal
		for i := 0; i <= int(SPAWN_BOUNDARY)*2; i += density {
			r, err := sf.NewRectangleShape()
			if err != nil {
				panic(err)
			}
			r.SetSize(v.NewCartesian(SPAWN_BOUNDARY*2, width).ToSFMLVector2f())
			r.SetPosition(worldToFramePosition(FRAME, v.NewCartesian(0, float64(i))).Sub(v.NewCartesian(SPAWN_BOUNDARY, SPAWN_BOUNDARY)).ToSFMLVector2f())
			r.SetOutlineThickness(0)
			r.SetFillColor(c)
			WINDOW.Draw(r, sf.DefaultRenderStates())
		}
	}(125, 2.0)

	// render player boundary
	func(boundaryExtension float64, thickness float64) {
		boundaryCenterPosition := worldToFramePosition(FRAME, v.NewZeroVector())
		rp := o.RenderProperties_s{
			thickness,
			0,
			p.TRANSPARENT,
			p.WHITE,
			v.NewUnitVector(),
		}
		c := CreateCircle(PLAYER_BOUNDARY+boundaryExtension, rp)
		c.SetPosition(boundaryCenterPosition.ToSFMLVector2f())
		WINDOW.Draw(c, sf.DefaultRenderStates())
	}(10, 5)

	// render player
	func() {
		framedPlayerPosition := worldToFramePosition(FRAME, PLAYER.Position)
		playerCircleShape := o.GetCircleShape(PLAYER.Circle_s, PLAYER.RenderProperties_s)
		playerCircleShape.SetPosition(framedPlayerPosition.ToSFMLVector2f())
		playerCircleShape.SetFillColor(p.RandomColor())

		// handle player bump
		func() {
			if PLAYER_BUMP_TIMEOUT.IsInBumpTimeout {
				secondsSinceLastBump := time.Since(PLAYER_BUMP_TIMEOUT.LastBumpTime).Seconds()
				segment := int(math.Floor(secondsSinceLastBump / BLINK_DURATION))
				c := p.STOP_SIGN_RED
				if segment%2 == 0 {
					c.A = 255
				} else {
					c.A = 0
				}
				playerCircleShape.SetFillColor(c)
			}
		}()

		WINDOW.Draw(playerCircleShape, sf.DefaultRenderStates())
	}()

	// render player control radius
	func() {
		framedPlayerPosition := worldToFramePosition(FRAME, PLAYER.Position)
		rp := o.RenderProperties_s{
			5,
			0,
			p.TRANSPARENT,
			func() sf.Color { w := p.WHITE; w.A = 30; return w }(),
			v.NewCartesian(1, 1),
		}
		c := CreateCircle(PLAYER_SETTINGS.ControlRadius, rp)
		c.SetPosition(framedPlayerPosition.ToSFMLVector2f())
		WINDOW.Draw(c, sf.DefaultRenderStates())
	}()

	// render asteroids
	func() {
		for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
			asteroid := e.Value.(*o.Asteroid_s)
			framedAsteroidPosition := worldToFramePosition(FRAME, asteroid.Position)
			asteroidCircleShape := o.GetCircleShape(asteroid.Circle_s, asteroid.RenderProperties_s)
			asteroidCircleShape.SetPosition(framedAsteroidPosition.ToSFMLVector2f())
			if LIGHT_SPAWN_DURATION == 02471 {
				asteroidCircleShape.SetFillColor(p.RandomColor())
			}
			WINDOW.Draw(asteroidCircleShape, sf.DefaultRenderStates())
		}
	}()

	// render slips
	func() {
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
				rp := o.RenderProperties_s{
					0,
					slip.End1.Sub(slip.End2).GetAngle(),
					color,
					p.WHITE,
					v.NewCartesian(1, 1),
				}
				r := CreateRectangle(v.NewCartesian(float64(SPAWN_BOUNDARY)*2.3, slip.Width), rp)
				r.SetPosition(worldToFramePosition(FRAME, slip.Segment_s.GetMidpoint()).ToSFMLVector2f())

				WINDOW.Draw(r, sf.DefaultRenderStates())
			}
		}
	}()

	// render health bar
	func() {
		healthPercentage := HEALTH_CURRENT / HEALTH_MAX
		red := p.RandomColor()
		red.A = 200
		rp := o.RenderProperties_s{
			0,
			0,
			red,
			p.WHITE,
			v.NewUnitVector(),
		}
		width := 15.0
		r := CreateRectangle(v.NewCartesian(healthPercentage*float64(GRAPHICS_SETTINGS.WindowX), width), rp)
		r.SetPosition(sf.Vector2f{float32(GRAPHICS_SETTINGS.WindowX) / 2.0, float32(width) / 2.0})
		WINDOW.Draw(r, sf.DefaultRenderStates())
	}()

	// // render boost bar
	// r, err := sf.NewRectangleShape()
	// if err != nil {
	// 	panic(err)
	// }
	// color := p.RED
	// color.A = 200
	// boostRatio := CURRENT_BOOST / MAX_BOOST
	// r.SetPosition(v.NewZeroVector().ToSFMLVector2f())
	// r.SetSize(v.NewCartesian(boostRatio*float64(GRAPHICS_SETTINGS.WindowX), 10).ToSFMLVector2f())
	// r.SetRotation(0)
	// r.SetFillColor(color)
	// r.SetOutlineThickness(0)
	// WINDOW.Draw(r, sf.DefaultRenderStates())

	WINDOW.Display()
}
