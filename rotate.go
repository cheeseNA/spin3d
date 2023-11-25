package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

type Point struct {
	X, Y, Z float64
}

func (p *Point) RotateX(angle float64) Point {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	y := p.Y*cos - p.Z*sin
	z := p.Y*sin + p.Z*cos
	return Point{p.X, y, z}
}

func (p *Point) RotateY(angle float64) Point {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	x := p.X*cos - p.Z*sin
	z := p.X*sin + p.Z*cos
	return Point{x, p.Y, z}
}

func (p *Point) RotateZ(angle float64) Point {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	x := p.X*cos - p.Y*sin
	y := p.X*sin + p.Y*cos
	return Point{x, y, p.Z}
}

func (p *Point) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

type Particle struct {
	Position   Point
	Luminosity float64
}

type Object interface {
	Particles(lightvec *Point) []Particle
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

func (d *Donut) Particles(lightvec *Point) []Particle {
	var particles []Particle
	for theta := 0.0; theta < 2.0*math.Pi; theta += d.dtheta {
		for phi := 0.0; phi < 2.0*math.Pi; phi += d.dphi {
			point := Point{
				X: d.r2 + d.r1*math.Cos(theta),
				Y: d.r1 * math.Sin(theta),
				Z: 0.0,
			}
			point = point.RotateY(phi)
			point = point.RotateX(d.angle1)
			point = point.RotateZ(d.angle2)

			normal := Point{
				X: math.Cos(theta),
				Y: math.Sin(theta),
				Z: 0.0,
			}
			normal = normal.RotateY(phi)
			normal = normal.RotateX(d.angle1)
			normal = normal.RotateZ(d.angle2)
			particles = append(particles, Particle{
				Position:   point,
				Luminosity: math.Max(0.0, normal.X*lightvec.X+normal.Y*lightvec.Y+normal.Z*lightvec.Z) / normal.Length(),
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
		//if particle.Luminosity <= 0.0 {
		//	continue
		//}

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
		dtheta: 0.02,
		dphi:   0.02,
		r1:     100.0,
		r2:     200.0,
		angle1: 0.0,
		angle2: 0.0,
		speed1: 0.05,
		speed2: 0.03,
	}

	donut.Tick()

	screen := &Screen{
		Width:  300,
		Height: 60,
		// k1: screen_width*K2*3/(8*(R1+R2));
		K1:    80 * 500 * 3 / (10 * (300.0)),
		K2:    400.0,
		Chars: ".,-~:;=!*#$@",
	}

	luminance := Point{0.0, 1.0, -1.0}

	screen.Init()

	interval := 100 * time.Millisecond

	stop := false

	go func() {
		for {
			time.Sleep(interval)
			if stop {
				continue
			}
			donut.Tick()
			screen.Clear()
			screen.Project(donut.Particles(&luminance))
			fmt.Print("\033[H")
			screen.Draw()
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			args := strings.Split(input, " ")
			switch args[0] {
			case "s1":
				donut.speed1 = 0.0
				_, err := fmt.Sscanf(args[1], "%f", &donut.speed1)
				if err != nil {
					fmt.Println(err)
				}
			case "s2":
				donut.speed2 = 0.0
				_, err := fmt.Sscanf(args[1], "%f", &donut.speed2)
				if err != nil {
					fmt.Println(err)
				}
			case "k1":
				screen.K1 = 0.0
				_, err := fmt.Sscanf(args[1], "%f", &screen.K1)
				if err != nil {
					fmt.Println(err)
				}
			case "k2":
				screen.K2 = 0.0
				_, err := fmt.Sscanf(args[1], "%f", &screen.K2)
				if err != nil {
					fmt.Println(err)
				}
			case "l":
				luminance.X = 0.0
				luminance.Y = 0.0
				luminance.Z = 0.0
				_, err := fmt.Sscanf(args[1], "%f", &luminance.X)
				if err != nil {
					fmt.Println(err)
				}
				_, err = fmt.Sscanf(args[2], "%f", &luminance.Y)
				if err != nil {
					fmt.Println(err)
				}
				_, err = fmt.Sscanf(args[3], "%f", &luminance.Z)
				if err != nil {
					fmt.Println(err)
				}
			case "stop":
				stop = true
			case "start":
				stop = false
			}
		}
	}()

	select {}
}
