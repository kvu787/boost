What we are trying to express

- Polygons on a 2D plane
- Position is modified through velocity and accleration
- Visual representation = data representation

Avoid the asteroids.

Input

- angle
- clicked?

Player

- position
- velocity
- acceleration

Asteroid

- velocity
- mass

Interactions

- player x asteroid => dead
- [stretch]: elastic collisions between asteroids

Invariants

- Player never starts too close to an asteroid


Installing the GoSFML2 library 
==============================

Note: the following is only tested for Windows, 64-bit

GoSFML2 is a media library used to render the graphics and GUI.

- Install TDM-GCC from http://tdm-gcc.tdragon.net/ with the Windows installer
- Add TDM-GCC to your PATH environment variable (should be at `C:\TDM-GCC-64\bin`)
- Download CSFML: http://www.sfml-dev.org/download/csfml/
	- Extract the zip
	- Move the `CSFML-2.1` folder the drive root
- Downlad GoSFML2 with `go get -u -d bitbucket.org/krepa098/gosfml2`
- In `bitbucket.org/krepa098/gosfml2/cgo.go`:
	- Add a line after the package statement and before `import "C"` with this format:
		- `// #cgo CFLAGS: -I<absolute path to include folder>`
		- Example: `// #cgo CFLAGS: -I/CSFML-2.1/include -I/CSFML-2.1/lib/gcc`
	- Change `LDFLAGS` line according to this format:
		- `// #cgo LDFLAGS: -L<absolute path to gcc libs> -lcsfml-window -lcsfml-graphics -lcsfml-audio`
		- Example: `// #cgo LDFLAGS: -L/CSFML-2.1/lib/gcc -lcsfml-window -lcsfml-graphics -lcsfml-audio`
- Run `go install` in the project root (`bitbucket.org/krepa098/gosfml2`)

```
func update() {

	// // spawn outside the spawn boundary
	// if ASTEROID_COUNT < ASTEROID_LIMIT {
	// 	i := rand.Intn(100)
	// 	if i == 0 {
	// 		ASTEROID_COUNT++
	// 		newAsteroid := func() Tagged {
	// 			radius := (rand.Float32() * 22.0) + 40.0
	// 			color := palette.BLUE
	// 			position := NewPolar(rand.Float64()*float64(SPAWN_BOUNDARY-5), NewRadians(2.0*math.Pi*rand.Float64()))
	// 			velocity := NewPolar((rand.Float64()*20.0)+20.0, NewRadians(2.0*math.Pi*rand.Float64()))
	// 			return &asteroid_s{
	// 				transform_s{position, velocity, NewZeroVector()},
	// 				circle_s{radius, 0, 0, color, palette.WHITE},
	// 				false,
	// 			}
	// 		}()
	// 		GAME_OBJECTS.PushFront(newAsteroid)
	// 	}
	// }

	// var camera *camera_s = listWhere(GAME_OBJECTS, CameraTag).(*camera_s)
	// var player *player_s = listWhere(GAME_OBJECTS, PlayerTag).(*player_s)
	// var asteroidsList *list.List = listSelect(GAME_OBJECTS, AsteroidTag)

	// // if mouse held, apply acceleration to player
	// if INPUT.isMousePressed {
	// 	playerFramePosition := getFramePosition(WINDOW_SIZE_X, WINDOW_SIZE_Y, camera, player.transform.position)
	// 	player.transform.acceleration = playerFramePosition.Sub(INPUT.mousePosition)
	// } else {
	// 	player.transform.acceleration = NewZeroVector()
	// }

	// // update player transform
	// player.transform = player.transform.applyAcceleration(secondsPerFrame)

	// // check if player is out of bounds
	// if player.transform.position.GetMagnitude() >= float64(PLAYER_BOUNDARY) {
	// 	fmt.Println("out of bounds foo")
	// 	// move back forcibly
	// 	player.transform.position = NewPolar(float64(PLAYER_BOUNDARY), player.transform.position.GetAngle())
	// 	// complete stop
	// 	player.transform.velocity = NewZeroVector()
	// 	player.transform.acceleration = NewZeroVector()
	// }

	// // update asteroid transforms
	// for e := asteroidsList.Front(); e != nil; e = e.Next() {
	// 	asteroid := e.Value.(*asteroid_s)
	// 	// check if asteroid is out of bounds
	// 	if asteroid.transform.position.GetMagnitude() >= float64(SPAWN_BOUNDARY) {
	// 		asteroid.shouldRemove = true
	// 	} else {
	// 		asteroid.transform = asteroid.transform.applyAcceleration(secondsPerFrame)
	// 	}
	// }

	// // remove dead asteroids
	// for e, next := GAME_OBJECTS.Front(), new(list.Element); e != nil; e = next {
	// 	next = e.Next()
	// 	if asteroid, ok := e.Value.(*asteroid_s); ok {
	// 		if asteroid.shouldRemove {
	// 			ASTEROID_COUNT--
	// 			GAME_OBJECTS.Remove(e)
	// 		}
	// 	}
	// }

	// // collide asteroids
	// for e1 := asteroidsList.Front(); e1 != nil; e1 = e1.Next() {
	// 	for e2 := e1.Next(); e2 != nil; e2 = e2.Next() {
	// 		a1 := e1.Value.(*asteroid_s)
	// 		a2 := e2.Value.(*asteroid_s)

	// 		distance := a1.transform.position.Sub(a2.transform.position).GetMagnitude()
	// 		sum_radius := a1.circle_s.radius + a2.circle_s.radius
	// 		is_intersecting := distance < float64(sum_radius)
	// 		if is_intersecting {
	// 			// wimpy collision resolution
	// 			resolve := func(x1, x2 transform_s) Vector {
	// 				fromx1tox2 := x2.position.Sub(x1.position)
	// 				proj := x1.velocity.Projection(fromx1tox2)
	// 				rej := x1.velocity.Rejection(fromx1tox2)
	// 				proj = proj.Mul(-1)
	// 				return proj.Add(rej)
	// 			}
	// 			a1.transform.velocity = resolve(a1.transform, a2.transform)
	// 			a2.transform.velocity = resolve(a2.transform, a1.transform)
	// 		}
	// 	}
	// }

	// // update camera
	// camera.Vector = player.transform.position
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

```