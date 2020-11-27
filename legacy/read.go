package main

import (
        "bufio"
        "fmt"
        "os"
        "strings"
)

func main() {
        type name struct {
            fname string
            lname string
        }
        var slicer []name
        var name_list []string
        var filename string
        fmt.Printf("Enter file name: \n")
        fmt.Scan(&filename)
        file, err := os.Open(filename)
        if err != nil {
                fmt.Println(err)
        }
        defer file.Close()
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                name_list = append(name_list, scanner.Text())
        }

        for _, each_name := range name_list {
                split_names := strings.Split(each_name, " ")
                person := name{}
                person.fname = split_names[0]
                person.lname = split_names[1]
                slicer = append(slicer, person)
        }
        for _, person_name := range slicer {
                output_name := string(person_name.fname) + " " + string(person_name.lname)
                fmt.Println(output_name)        
        
        }
}
