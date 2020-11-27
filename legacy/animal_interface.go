package main

import (
	"fmt"
)

type Animal interface {
	Eat()
	Move()
	Speak()
	Name() string
}

type Cow struct {
	name string
}

func (x Cow) Eat() {
	fmt.Println("grass")
}

func (x Cow) Move() {
	fmt.Println("walk")
}

func (x Cow) Speak() {
	fmt.Println("moo")
}

func (x Cow) Name() string {
	return x.name
}

type Bird struct {
	name string
}

func (x Bird) Eat() {
	fmt.Println("worms")
}

func (x Bird) Move() {
	fmt.Println("fly")
}

func (x Bird) Speak() {
	fmt.Println("peep")
}

func (x Bird) Name() string {
	return x.name
}

type Snake struct {
	name string
}

func (x Snake) Eat() {
	fmt.Println("mice")
}

func (x Snake) Move() {
	fmt.Println("slither")
}

func (x Snake) Speak() {
	fmt.Println("hss")
}

func (x Snake) Name() string {
	return x.name
}

func Eat(r Animal) {
	r.Eat()
}

func Move(r Animal) {
	r.Move()
}

func Speak(r Animal) {
	r.Speak()
}

func Name(r Animal) string {
	return r.Name()
}

func main() {
	var animal_list []Animal
	fmt.Println("Create a new animal or query an existing animal:")
	fmt.Println("Example: \"newanimal <name> <animal>\" or \"query <name> <action>\"")
	fmt.Println("<Enter * to exit the program> \n")
	for {
		var command, name, animal_or_action string
		fmt.Printf("> ")
		fmt.Scan(&command)
		if command == "*" {
			break
		}
		fmt.Scan(&name)
		fmt.Scan(&animal_or_action)
		switch command {
		case "newanimal":
			switch animal_or_action {
			case "cow":
				animal_list = append(animal_list, Cow{name: name})
			case "bird":
				animal_list = append(animal_list, Bird{name: name})
			case "snake":
				animal_list = append(animal_list, Snake{name: name})
			default:
				fmt.Println("Pick an animal between cow, bird and snake")
			}
		case "query":
			for _, animal := range animal_list {
				animal_name := Name(animal)
				if animal_name == name {
					switch animal_or_action {
					case "eat":
						Eat(animal)
					case "move":
						Move(animal)
					case "speak":
						Speak(animal)
					default:
						fmt.Println("Pick an action between eat, move and speak")
					}
				} else {
					fmt.Println("An animal with that name does not exist yet")
				}

			}
		}
	}
}
