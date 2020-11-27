package main

import (
        "encoding/json"
        "fmt"
        "os"
)

func main() {
        person := make(map[string]string)
        var name_input string
        var address_input string

        fmt.Printf("Enter name: \n")
        fmt.Scan(&name_input)
        fmt.Printf("Enter address: \n")
        fmt.Scan(&address_input)

        person["name"] = name_input
        person["address"] = address_input

        barr, e := json.Marshal(person)
        if e != nil {
                fmt.Printf("Unable to convert map to json object: \n %s\n", e)
                os.Exit(1)
        }
        fmt.Println(string(barr))

}