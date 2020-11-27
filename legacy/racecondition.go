package main

import (
        "fmt"
        "sync"
        "time"
)

//Here a race condition is created by first initialising a count integer
var count int

// a named function to be called as a goroutine is created
func routine(wg *sync.WaitGroup) {
        defer wg.Done()
        for x := 0; x < 50; x++ {
                count++
        }
        fmt.Printf("goroutine named function call running \n")
        //a delay is introduced to showcase the race condition
        time.Sleep(time.Duration(2) * time.Second)
}

func main() {
        //a wait group is created to be used for the named function call
        var wg sync.WaitGroup
        wg.Add(1)
        //the named function is called here
        go routine(&wg)
        //an anonymous function reading and writing to the count variable is called as a goroutine here
        go func() {
                count = count + 69
                //due to the delay in the named function, the count value below will be 69
                fmt.Println("race condition count: ",count)
                fmt.Printf("goroutine anon function call running \n")
        }()
        wg.Wait()
        //after the named function goroutine is complete, the final count shows a different value due to overlapped use of the shared variable "count". This results case of a race condition
        fmt.Println("final count:",count)
}
