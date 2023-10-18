package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
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
