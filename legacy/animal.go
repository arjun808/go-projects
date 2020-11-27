package main

import (
	"fmt"
)

type Animal struct {
	eat   string
	move  string
	speak string
}

func (x Animal) Eat() {
	fmt.Println(x.eat)
}

func (x Animal) Move() {
	fmt.Println(x.move)
}

func (x Animal) Speak() {
	fmt.Println(x.speak)
}

func main() {
	cow := Animal{"grass", "walk", "moo"}
	bird := Animal{"worms", "fly", "peep"}
	snake := Animal{"mice", "slither", "hss"}

	fmt.Println("Enter an animal name and and action: e.g. \"snake eat\"")
	fmt.Println("<Enter * to exit the program>")
	for {
		var name, action string
		fmt.Printf("> ")
		fmt.Scan(&name)
		if name == "*" {
			break
		}
		fmt.Scan(&action)
		var animal Animal
		switch name {
		case "cow":
			animal = cow
		case "bird":
			animal = bird
		case "snake":
			animal = snake
		}

		switch action {
		case "eat":
			animal.Eat()
		case "move":
			animal.Move()
		case "speak":
			animal.Speak()
		}
	}
}
