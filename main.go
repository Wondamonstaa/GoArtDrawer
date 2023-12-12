package main

/*
Name: Kiryl Baravikou
UIN: 656339218
Project 3: Go Interfaces
Date: 11/30/2023
Professor: Jon Solworth
Description:
The goal of the following project is to create a program that will draw the specified shapes on the screen.
The colors are specified in the colors' matrix. The program will be able to draw the following shapes:
	1. Rectangle
	2. Circle
	3. Triangle
The program will be able to draw the shapes on the screen and save the output as a .ppm file.
*/

import (
	"fmt"
)

func main() {
	fmt.Println("starting ...")
	display.initialize(1024, 1024)

	rect := Rectangle{Point{100, 300}, Point{600, 900}, RED}
	err := rect.draw(&display)
	if err != nil {
		fmt.Println("rect: ", err)
	}

	rect2 := Rectangle{Point{0, 0}, Point{100, 1024}, GREEN}
	err = rect2.draw(&display)
	if err != nil {
		fmt.Println("rect2: ", err)
	}

	rect3 := Rectangle{Point{0, 0}, Point{100, 1022}, 102}
	err = rect3.draw(&display)
	if err != nil {
		fmt.Println("rect3: ", err)
	}

	circ := Circle{Point{500, 500}, 200, GREEN}
	err = circ.draw(&display)
	if err != nil {
		fmt.Println("circ: ", err)
	}

	tri := Triangle{Point{100, 100}, Point{600, 300}, Point{859, 850}, YELLOW}
	err = tri.draw(&display)
	if err != nil {
		fmt.Println("tri: ", err)
	}

	//Will display the output on the screen in the ppm format
	display.screenShot("output")
}
