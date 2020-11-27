package main

import (
        "bufio"
        "fmt"
        "math"
        "os"
        "sort"
        "strconv"
        "strings"
        "sync"
)

func main() {
mainloop:
        for {
                fmt.Printf("Enter integers separated by space \n")
                fmt.Printf("> ")
                r := bufio.NewReader(os.Stdin)
                input_list, _ := r.ReadString('\n')
                s := strings.Split(strings.TrimSpace(input_list), " ")
                var arr []int

                for _, str := range s {
                        num, err := strconv.Atoi(str)
                        if err != nil {
                                fmt.Println("Please enter only integers \n")
                                continue mainloop
                        }
                        arr = append(arr, num)
                }
                fmt.Println(arr)
                parts := 4

                var wg sync.WaitGroup
                c := make(chan []int, parts)
                n := int(math.Max(math.Ceil(float64(len(s))/float64(parts)), 1))
                for i := 0; i < parts; i++ {
                        start := int(math.Min(float64(i*n), float64(len(arr))))
                        end := int(math.Min(float64((i+1)*n), float64(len(arr))))

                        wg.Add(1)

                        go func(a []int) {
                                fmt.Println("sorting: ", a)
                                sort.Ints(a)
                                c <- a
                                wg.Done()
                        }(arr[start:end])
                }

                wg.Wait()
                close(c)
                var final_arr []int
                for item := range c {
                        for _, z := range item {
                                final_arr = append(final_arr, z)
                        }
                }
                sort.Ints(final_arr)
                fmt.Println("sorted: ", final_arr)
        }

}
