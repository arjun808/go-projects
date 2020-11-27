package main

import (
        "fmt"
        "strings"
)

func main() {
        var string_input string
        fmt.Printf("Enter a string: ")
        fmt.Scan(&string_input)
        checkstring := strings.ToLower(strings.TrimSpace(string_input))
        if strings.Contains(checkstring, "a") && checkstring[:1] == "i" && checkstring[len(checkstring)-1:] == "n" {
                fmt.Printf("Found! \n")
        } else {
                fmt.Printf("Not Found! \n")

        }

}