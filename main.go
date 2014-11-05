package main

import (
	"container/list"
	"fmt"
	"runtime"
	"time"

	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

// config variables
const (
	WINDOW_SIZE_X uint = 1000
	WINDOW_SIZE_Y uint = 720
	FPS           int  = 65
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

	// update asteroid transforms
	for e := asteroidsList.Front(); e != nil; e = e.Next() {
		asteroid := e.Value.(*asteroid_s)
		asteroid.transform = asteroid.transform.applyAcceleration(secondsPerFrame)
	}

	// // collide asteroids
	// for e1 := asteroidsList.Front(); e1 != nil; e1 = e1.Next() {
	// 	for e2 := asteroidsList.Front(); e2 != nil; e2 = e2.Next() {
	// 	}
	// }

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
