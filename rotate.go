package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type Point struct {
	X, Y, Z float64
}

type Vector struct {
	X, Y, Z float64
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

type Particle struct {
	Position   Point
	Luminosity float64
}

type Object interface {
	Particles(lightvec Vector) []Particle
	Tick()
}

type Donut struct {
	dtheta float64
	dphi   float64
	r1     float64
	r2     float64
	angle1 float64
	angle2 float64
	speed1 float64
	speed2 float64
}

func (d *Donut) Particles(lightvec Vector) []Particle {
	var particles []Particle
	for theta := 0.0; theta < 2.0*math.Pi; theta += d.dtheta {
		for phi := 0.0; phi < 2.0*math.Pi; phi += d.dphi {
			ox := (d.r2 + d.r1*math.Cos(theta)) * math.Cos(phi)
			oy := d.r1 * math.Sin(theta)
			oz := (d.r2 + d.r1*math.Cos(theta)) * math.Sin(phi)

			y := oy*math.Cos(d.angle1) + oz*math.Sin(d.angle1)
			z := -oy*math.Sin(d.angle1) + oz*math.Cos(d.angle1)

			oy = y
			oz = z

			x := ox*math.Cos(d.angle2) + oz*math.Sin(d.angle2)
			y = -ox*math.Sin(d.angle2) + oz*math.Cos(d.angle2)
			z = oz
			particles = append(particles, Particle{
				Position:   Point{x, y, z},
				Luminosity: 1.0,
			})
		}
	}
	return particles
}

func (d *Donut) Tick() {
	d.angle1 += d.speed1
	d.angle2 += d.speed2
}

type Screen struct {
	Width     int
	Height    int
	K1        float64
	K2        float64
	Chars     string
	zbuffer   [][]float64
	luminance [][]float64
}

func (s *Screen) Init() {
	s.zbuffer = make([][]float64, s.Width)
	s.luminance = make([][]float64, s.Width)
	for i := range s.zbuffer {
		s.zbuffer[i] = make([]float64, s.Height)
		s.luminance[i] = make([]float64, s.Height)
	}
}

func (s *Screen) Clear() {
	for i := range s.zbuffer {
		for j := range s.zbuffer[i] {
			s.zbuffer[i][j] = 0.0
			s.luminance[i][j] = 0.0
		}
	}
}

func (s *Screen) Project(particles []Particle) {
	for _, particle := range particles {
		ooz := 1.0 / (particle.Position.Z + s.K2)
		x := particle.Position.X * ooz * s.K1
		y := particle.Position.Y * ooz * s.K1
		sx := int(x) + s.Width/2
		sy := int(y/2) + s.Height/2
		if sx < 0 || sx >= s.Width || sy < 0 || sy >= s.Height {
			continue
		}
		if particle.Luminosity <= 0.0 {
			continue
		}

		// bigger ooz means closer to the screen
		if ooz > s.zbuffer[sx][sy] {
			s.zbuffer[sx][sy] = ooz
			s.luminance[sx][sy] = particle.Luminosity
		}
	}
}

func (s *Screen) Draw() {
	for j := 0; j < s.Height; j++ {
		for i := 0; i < s.Width; i++ {
			luminance := s.luminance[i][j]
			if luminance <= 0.0 {
				fmt.Print(" ")
			} else {
				index := int(luminance * float64(len(s.Chars)))
				if index >= len(s.Chars) {
					index = len(s.Chars) - 1
				}
				fmt.Print(string(s.Chars[index]))
			}
		}
		fmt.Println()
	}
}

func main() {
	donut := &Donut{
		dtheta: 0.07,
		dphi:   0.02,
		r1:     1.0,
		r2:     2.0,
		angle1: 0.0,
		angle2: 0.0,
		speed1: 0.05,
		speed2: 0.03,
	}

	donut.Tick()

	screen := &Screen{
		Width:  80,
		Height: 40,
		K1:     80 * 5 * 3 / (8 * (3.0)),
		K2:     5.0,
		Chars:  ".,-~:;=!*#$@",
	}

	screen.Init()

	interval := 50 * time.Millisecond

	go func() {
		for {
			time.Sleep(interval)
			donut.Tick()
			screen.Clear()
			screen.Project(donut.Particles(Vector{0.0, 0.0, 1.0}))
			fmt.Print("\033[H\033[2J")
			screen.Draw()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			// accept int and float
			floatInput, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input.")
				continue
			}
			fmt.Println("Input:", floatInput)
		}
	}()

	select {}
}
