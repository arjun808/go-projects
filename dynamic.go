package main

import "fmt"

func main() {
    subnets := []string{"web", "db", "bread", "bake", "fake", "qwerty", "asdfg", "ecdcx", "aasdw"}
    az := []string{"a","b","c"}
    l := len(az)
    recurse := 0
   for i,_ := range subnets {
       if recurse < l {
      fmt.Println("az1: ", az[recurse], "subnets: ", subnets[i])
      recurse++
     }else {
       recurse=0
}
}
}
