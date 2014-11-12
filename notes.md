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
