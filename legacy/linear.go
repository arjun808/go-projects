package main

import (
        "fmt"
)

func GenDisplaceFn(a, v0, s0 float64) func(float64) float64 {
        displacement := func(t float64) float64 {
                return 0.5*a*t*t + v0*t + s0
        }
        return displacement
}

func main() {
        var acceleration, init_velocity, init_displacement, time float64
        fmt.Printf("Enter acceleration value:\n")
        fmt.Scan(&acceleration)
        fmt.Printf("Enter initial velocity value:\n")
        fmt.Scan(&init_velocity)
        fmt.Printf("Enter initial displacement value:\n")
        fmt.Scan(&init_displacement)
        fn := GenDisplaceFn(acceleration, init_velocity, init_displacement)
        fmt.Printf("Enter time value:\n")
        fmt.Scan(&time)
        fmt.Printf("\nThe displacement output is: ")
        fmt.Println(fn(time))
}
