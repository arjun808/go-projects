package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "strings"
  b64 "encoding/base64"
  "strconv"
)

func kv_get(AppVersion string) (string) {
var response string
kv_get_resp, err := http.Get("http://172.31.4.66:8500/v1/kv/"+AppVersion)
if err != nil {
       fmt.Println(err)
}else {
  body, _ := ioutil.ReadAll(kv_get_resp.Body)
  text := string(body)
  if text != "" {
  s := strings.Split(text, "\"")
  x := s[11]
  data, _ := b64.StdEncoding.DecodeString(x)
  if data[0] == '"' {
        data = data[1:]
    }
  if i := len(data)-1; data[i] == '"' {
        data = data[:i]
    }
  response = string(data)
   }else {
	response = ""
      }
    defer kv_get_resp.Body.Close()
     }
	 return response
}

func main() {
AppV := "1.0"
data := kv_get(AppV)
if data != "" {
 r := "-"
 x := strings.LastIndex(data, r) + 1
 Nodenumber, _ := strconv.ParseInt(data[x:], 0, 64)
 Nodenumber++ 
 fmt.Println(Nodenumber)
}
}
