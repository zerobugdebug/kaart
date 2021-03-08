package main

import "log"

func main() {
	total := 0
	for a := 0; a <= 12; a++ {
		for b := 0; b <= 12-a; b++ {
			for c := 0; c <= 12-a-b; c++ {
				d := 12 - a - b - c
				log.Print(a, b, c, d)
				total++
			}
		}
	}
	log.Print(total)
}
