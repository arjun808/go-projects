package main

import (
  "fmt"
  "net/http"
)

func main() {

req, err := http.NewRequest("PUT", "http://172.31.3.82:8500/v1/agent/service/maintenance/apache2?enable=false", nil)
//req, err := http.NewRequest("PUT", "http://172.31.3.82:8500/v1/agent/service/maintenance/apache2?enable=true&reason=For+glory", nil)
if err != nil {
 fmt.Println(err)
}
resp, err := http.DefaultClient.Do(req)
if err != nil {
 fmt.Println(err)
}
defer resp.Body.Close()
}

