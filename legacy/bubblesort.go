package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func BubbleSort(slice []int) {
	incomplete := true
	for incomplete {
		incomplete = false
		for iter := 1; iter < len(slice); iter++ {
			if slice[iter] < slice[iter-1] {
				Swap(slice, iter)
				incomplete = true
			}
		}
	}
}

func Swap(slice []int, i int) {
	slice[i-1], slice[i] = slice[i], slice[i-1]
}

func main() {
	var input_string string
	var input_slice []int
	fmt.Printf("Enter an integer sequence separated by commas:\n")
	fmt.Scan(&input_string)
	input_string_slice := strings.Split(input_string, ",")
	for _, item := range input_string_slice {
		int_item, err := strconv.Atoi(item)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		input_slice = append(input_slice, int_item)
	}
	BubbleSort(input_slice)
	fmt.Println(input_slice)
}
