package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "strings"
  b64 "encoding/base64"
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

func kv_put(Nodename string, AppVersion string) (string) {
  var response string
  body := strings.NewReader(Nodename)
  req, req_err := http.NewRequest("PUT", "http://172.31.4.66:8500/v1/kv/"+AppVersion, body)
  if req_err != nil {
	response = "Request error"
    }else {
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, resp_err := http.DefaultClient.Do(req)
    if resp_err != nil {
	 response = "Response error"
    }else {
	   response = "true"
	   
}
defer resp.Body.Close()
}
  return response
}


func main() {
Nodename := "App1"
AppV := "1.0"
data := kv_get(AppV)

NewNodename := data + "," + Nodename

bool := kv_put(NewNodename, AppV)

if bool == "true" {
fmt.Println("done da")
}else {
fmt.Println("no da")
}

}
