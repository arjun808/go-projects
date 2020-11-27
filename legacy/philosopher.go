package main

import (
        "fmt"
        "sync"
)

func main() {
        host := make(chan int, 2)
        for x := 0; x < 5; x++ {
                ChopstickPool.Put(new(Chopstick))
        }

        phil := make([]*Philosopher, 5)
        for i := 0; i < 5; i++ {
                phil[i] = &Philosopher{i + 1}
        }

        for i := 0; i < 5; i++ {
                for j := 1; j < 4; j++ {
                        wg.Add(1)
                        go phil[i].Eat(host, j)
                }
        }
        host <- 1
        host <- 1
        wg.Wait()
}


type Chopstick struct {
}

type Philosopher struct {
        n int
}

var ChopstickPool = sync.Pool{
        New: func() interface{} {
                return new(Chopstick)
        },
}

var wg sync.WaitGroup

func (p Philosopher) Eat(host chan int, times int) {
        defer wg.Done()

        <-host

        fmt.Printf("Philosopher %d is eating \n", p.n)
        right := ChopstickPool.Get()
        left := ChopstickPool.Get()
        ChopstickPool.Put(right)
        ChopstickPool.Put(left)
        fmt.Printf("Philosopher %d is done eating %d \n", p.n, times)

        host <- 1
}

