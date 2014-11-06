package main

import (
	"fmt"
	"runtime"
	"time"

	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"
	. "bitbucket.org/kvu787/boost/objects"

	sf "bitbucket.org/kvu787/gosfml2"
)

// config variables

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

		input()  // update global INPUT
		update() // update global GAME_OBJECTS
		render() // write to screen

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

	// update asteroid transforms
	for e := ASTEROIDS.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*Asteroid_s)
		asteroid.Transform = asteroid.Transform.Act(DURATION_PER_FRAME)
	}
}

func render() {
	WINDOW.Clear(palette.BLACK)

	// get camera
	var camera Vector = PLAYER.Transform.Position.Add(CAMERA_OFFSET)

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
