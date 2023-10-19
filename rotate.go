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

type Particle struct {
	Position   Point
	Luminosity float64
}

type Object interface {
	Particles(lightvec Point) []Particle
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

func (d *Donut) Particles(lightvec Point) []Particle {
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

func main() {
	donut := &Donut{
		dtheta: 0.07,
		dphi:   0.02,
		r1:     1.0,
		r2:     2.0,
		angle1: 0.0,
		angle2: 0.0,
		speed1: 0.01,
		speed2: 0.01,
	}

	donut.Tick()

	interval := 1 * time.Second
	counter := 1.0
	step := 1.0

	go func() {
		for {
			time.Sleep(interval)
			fmt.Println(counter)
			counter += step
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
			step = floatInput
		}
	}()

	select {}
}
