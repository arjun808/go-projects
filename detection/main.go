package main

import (
 "fmt"
 "github.com/google/go-github/github"
 "context"
 "golang.org/x/oauth2"
 "log"
)

func main() {

ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "3f1b118578a19549a2d4432456163c683fd1d65b"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List(ctx, "", nil)

        if err != nil{
         fmt.Println("Error: ", err)
           }else {
             fmt.Println(repos)
               }

if _, ok := err.(*github.RateLimitError); ok {
	log.Println("hit rate limit")
}
}
