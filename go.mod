module rpi-rgb-screen

go 1.24.3

require github.com/mcuadros/go-rpi-rgb-led-matrix v0.0.0-20180401002551-b26063b3169a

require (
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20221208032759-85de2813cf6b // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20231223183121-56fa3ac82ce7 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/exp/shiny v0.0.0-20250531010427-b6e5de432a8b // indirect
	golang.org/x/image v0.27.0 // indirect
	golang.org/x/mobile v0.0.0-20250520180527-a1d90793fc63 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

replace github.com/mcuadros/go-rpi-rgb-led-matrix => ../go-rpi-rgb-led-matrix
