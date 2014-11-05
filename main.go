package main

import (
	"fmt"
	"runtime"
	"time"

	"bitbucket.org/kvu787/boost/lib/palette"
	. "bitbucket.org/kvu787/boost/lib/vector"

	sf "bitbucket.org/kvu787/gosfml2"
)

// config variables
const (
	WINDOW_SIZE_X uint = 640
	WINDOW_SIZE_Y uint = 480
	FPS           int  = 120
)

func main() {
	// game variables
	camera := NewZeroVector()
	player := player_s{
		NewZeroVector(),
		NewZeroVector(),
		NewZeroVector(),
	}
	startTime := time.Now()
	secondsPerFrame := time.Duration(int64(time.Second) / int64(FPS))
	var isMousePressed = false
	circle := func() *sf.CircleShape {
		circle, err := sf.NewCircleShape()
		var innerCircleRadius float64 = 5
		var outlineThickness float64 = 0.8
		if err != nil {
			panic(err)
		}
		circle.SetFillColor(palette.BLUE)
		circle.SetOutlineColor(palette.WHITE)
		circle.SetOutlineThickness(float32(outlineThickness))
		circle.SetRadius(float32(innerCircleRadius))
		circle.SetOrigin(NewCartesian(innerCircleRadius, innerCircleRadius).ToSFMLVector2f())
		circle.SetPosition(NewZeroVector().ToSFMLVector2f())
		circle.SetRotation(0)
		return circle
	}()
	asteroidPosition := NewCartesian(50, -50)
	asteroid := func() *sf.CircleShape {
		circle, err := sf.NewCircleShape()
		var innerCircleRadius float64 = 20
		var outlineThickness float64 = 3
		if err != nil {
			panic(err)
		}
		circle.SetFillColor(palette.GRAY)
		circle.SetOutlineColor(palette.WHITE)
		circle.SetOutlineThickness(float32(outlineThickness))
		circle.SetRadius(float32(innerCircleRadius))
		circle.SetOrigin(NewCartesian(innerCircleRadius, innerCircleRadius).ToSFMLVector2f())
		circle.SetPosition(asteroidPosition.ToSFMLVector2f())
		circle.SetRotation(0)
		return circle
	}()

	// setup
	fmt.Println("Setting up window...")
	runtime.LockOSThread()
	window := sf.NewRenderWindow(
		sf.VideoMode{640, 480, 32},
		"boost",
		sf.StyleDefault,
		sf.DefaultContextSettings())

	// main loop
	for window.IsOpen() {
		getRenderPosition := func(winx, winy uint, camera, x Vector) Vector {
			frame := camera.Add(NewCartesian(-0.5*float64(winx), -0.5*float64(winy)))
			return x.Sub(frame)
		}

		// loop variables
		var mouseClick Vector = nil

		// vsync
		startTime = time.Now()

		// check for user input
		for event := window.PollEvent(); event != nil; event = window.PollEvent() {
			switch event.(type) {
			case sf.EventClosed:
				window.Close()
			case sf.EventMouseButtonPressed:
				isMousePressed = true
			case sf.EventMouseButtonReleased:
				isMousePressed = false
			}
		}
		if isMousePressed {
			mousePosition := func() Vector {
				position := sf.MouseGetPosition(window)
				return NewCartesian(float64(position.X), float64(position.Y))
			}()
			mouseClick = mousePosition
		}

		// set player object acceleration
		if mouseClick != nil {
			playerRenderPosition := getRenderPosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, player.position)
			player.acceleration = playerRenderPosition.Sub(mouseClick)
		} else {
			player.acceleration = NewZeroVector()
		}

		// update player
		velocityDelta := player.acceleration.Mul(secondsPerFrame.Seconds())
		player.velocity = player.velocity.Add(velocityDelta)
		positionDelta := player.velocity.Mul(secondsPerFrame.Seconds())
		player.position = player.position.Add(positionDelta)

		// update camera
		camera = player.position.Add(NewCartesian(-100, -100))

		// render
		circle.SetPosition(getRenderPosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, player.position).ToSFMLVector2f())
		asteroid.SetPosition(getRenderPosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, asteroidPosition).ToSFMLVector2f())
		window.Clear(sf.ColorBlack())
		window.Draw(circle, sf.DefaultRenderStates())
		window.Draw(asteroid, sf.DefaultRenderStates())
		window.Display()

		// sleep if frametime is short
		time.Sleep(time.Duration(int64(secondsPerFrame) - int64(time.Since(startTime))))
	}
}

type player_s struct {
	position     Vector
	velocity     Vector
	acceleration Vector
}
