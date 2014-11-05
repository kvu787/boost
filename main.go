package main

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	. "bitbucket.org/kvu787/boost/lib/angle"
	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

// config variables
const (
	WINDOW_SIZE_X   uint = 1000
	WINDOW_SIZE_Y   uint = 720
	FPS             int  = 65
	PLAYER_BOUNDARY uint = 400
	SPAWN_BOUNDARY  uint = 500
	ASTEROID_LIMIT  uint = 40
)

func main() {
	// timing variables
	startTime := time.Now()
	durationPerFrame := time.Duration(int64(time.Second) / int64(FPS))

	// setup
	fmt.Println("Setting up window...")
	setup()

	// main loop
	fmt.Println("Entering game loop...")
	for WINDOW.IsOpen() {
		// vsync
		startTime = time.Now()

		input()                            // update global INPUT
		update(durationPerFrame.Seconds()) // update global GAME_OBJECTS
		render()                           // write to screen

		// sleep if frametime is short
		time.Sleep(time.Duration(int64(durationPerFrame) - int64(time.Since(startTime))))
	}
}

func getFramePosition(winx, winy uint, camera, x Vector) Vector {
	frame := camera.Add(NewCartesian(-0.5*float64(winx), -0.5*float64(winy)))
	return x.Sub(frame)
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
			INPUT.isMousePressed = true
		case sf.EventMouseButtonReleased:
			INPUT.isMousePressed = false
		}
	}
	position := sf.MouseGetPosition(WINDOW)
	INPUT.mousePosition = NewCartesian(float64(position.X), float64(position.Y))
}

func update(secondsPerFrame float64) {
	fmt.Println(ASTEROID_COUNT)
	if ASTEROID_COUNT < ASTEROID_LIMIT {
		ASTEROID_COUNT++
		newAsteroid := func() Tagged {
			radius := (rand.Float32() * 22.0) + 40.0
			color := palette.BLUE
			position := NewPolar(rand.Float64()*float64(SPAWN_BOUNDARY-5), NewRadians(2.0*math.Pi*rand.Float64()))
			velocity := NewPolar((rand.Float64()*20.0)+20.0, NewRadians(2.0*math.Pi*rand.Float64()))
			return &asteroid_s{
				transform_s{position, velocity, NewZeroVector()},
				circle_s{radius, 0, 0, color, palette.WHITE},
				false,
			}
		}()
		GAME_OBJECTS.PushFront(newAsteroid)
	}

	var camera *camera_s = listWhere(GAME_OBJECTS, CameraTag).(*camera_s)
	var player *player_s = listWhere(GAME_OBJECTS, PlayerTag).(*player_s)
	var asteroidsList *list.List = listSelect(GAME_OBJECTS, AsteroidTag)

	// if mouse held, apply acceleration to player
	if INPUT.isMousePressed {
		playerFramePosition := getFramePosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, player.transform.position)
		player.transform.acceleration = playerFramePosition.Sub(INPUT.mousePosition)
	} else {
		player.transform.acceleration = NewZeroVector()
	}

	// update player transform
	player.transform = player.transform.applyAcceleration(secondsPerFrame)

	// check if player is out of bounds
	if player.transform.position.GetMagnitude() >= float64(PLAYER_BOUNDARY) {
		fmt.Println("out of bounds foo")
		// move back forcibly
		player.transform.position = NewPolar(float64(PLAYER_BOUNDARY), player.transform.position.GetAngle())
		// complete stop
		player.transform.velocity = NewZeroVector()
		player.transform.acceleration = NewZeroVector()
	}

	// update asteroid transforms
	for e := asteroidsList.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*asteroid_s)
		// check if asteroid is out of bounds
		if asteroid.transform.position.GetMagnitude() >= float64(SPAWN_BOUNDARY) {
			asteroid.shouldRemove = true
		} else {
			asteroid.transform = asteroid.transform.applyAcceleration(secondsPerFrame)
		}
	}

	// remove dead asteroids
	for e, next := GAME_OBJECTS.Front(), new(list.Element); e != nil; e = next {
		next = e.Next()
		if asteroid, ok := e.Value.(*asteroid_s); ok {
			if asteroid.shouldRemove {
				ASTEROID_COUNT--
				GAME_OBJECTS.Remove(e)
			}
		}
	}

	// collide asteroids
	for e1 := asteroidsList.Front(); e1 != nil; e1 = e1.Next() {
		for e2 := e1.Next(); e2 != nil; e2 = e2.Next() {
			a1 := e1.Value.(*asteroid_s)
			a2 := e2.Value.(*asteroid_s)

			distance := a1.transform.position.Sub(a2.transform.position).GetMagnitude()
			sum_radius := a1.circle_s.radius + a2.circle_s.radius
			is_intersecting := distance < float64(sum_radius)
			if is_intersecting {
				// wimpy collision resolution
				resolve := func(x1, x2 transform_s) Vector {
					fromx1tox2 := x2.position.Sub(x1.position)
					proj := x1.velocity.Projection(fromx1tox2)
					rej := x1.velocity.Rejection(fromx1tox2)
					proj = proj.Mul(-1)
					return proj.Add(rej)
				}
				a1.transform.velocity = resolve(a1.transform, a2.transform)
				a2.transform.velocity = resolve(a2.transform, a1.transform)
			}
		}
	}

	// update camera
	camera.Vector = player.transform.position
}

func render() {
	var camera *camera_s = listWhere(GAME_OBJECTS, CameraTag).(*camera_s)
	var player *player_s = listWhere(GAME_OBJECTS, PlayerTag).(*player_s)
	var asteroidsList *list.List = listSelect(GAME_OBJECTS, AsteroidTag)

	// render
	WINDOW.Clear(palette.BLACK)

	// render player
	playerDrawer := player.circle.GetDrawer()
	playerFramePosition := getFramePosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, player.transform.position).ToSFMLVector2f()
	playerDrawer.SetPosition(playerFramePosition)
	WINDOW.Draw(playerDrawer, sf.DefaultRenderStates())

	// render asteroids
	for e := asteroidsList.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*asteroid_s)
		asteroidDrawer := asteroid.circle_s.GetDrawer()
		asteroidFramePosition := getFramePosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, asteroid.transform.position).ToSFMLVector2f()
		asteroidDrawer.SetPosition(asteroidFramePosition)
		WINDOW.Draw(asteroidDrawer, sf.DefaultRenderStates())
	}

	WINDOW.Display()
}
