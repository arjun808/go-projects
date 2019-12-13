package main

import (
  "fmt"
  "os"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func main() {
tag := os.Args[1]
    resp, err := http.Get("http://172.31.1.198:8500/v1/catalog/service/apache2?tag="+tag)
     if err != nil {
        fmt.Println(err)
    }
      body, _ := ioutil.ReadAll(resp.Body)
      defer resp.Body.Close()

      byt := []byte(string(body))
      var data []map[string]interface{}
      var address []string
        if err := json.Unmarshal(byt, &data); err != nil {
            panic(err)
        }
        for i, _ := range data {
            x := data[i]["Address"].(string)
                address = append(address, x)
        }
        fmt.Println(address)
        for item, _ := range address {
        addr := address[item]
        req, err := http.NewRequest("PUT", "http://"+addr+":8500/v1/agent/leave", nil)
          if err != nil {
        fmt.Println(err)
        }
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
          fmt.Println(err)
       }
       defer resp.Body.Close()
        }
      }

