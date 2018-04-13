package main

import (
   "github.com/aws/aws-sdk-go/aws"
   "github.com/aws/aws-sdk-go/aws/session"
   "github.com/aws/aws-sdk-go/service/ec2"
   "net/http"
   "fmt"
   b64 "encoding/base64"
    	"encoding/json"
		"io/ioutil"
		  "strings"
		  "strconv"
		  "time"
		  
       )

func main(){
   http.HandleFunc("/",welcome)
   http.HandleFunc("/createins", Loop)
   http.ListenAndServe(":9090",nil)
}

func welcome(w http.ResponseWriter, r *http.Request) {

fmt.Fprintf(w,htmlStr)

}

type arraymap []map[string]interface{}

func maint(AppVersion string) {
fresh_tag := AppVersion
resp, err := http.Get("http://172.31.4.66:8500/v1/kv/?recurse=true")
if err != nil {
   fmt.Println(err)
}
defer resp.Body.Close()
body, _ := ioutil.ReadAll(resp.Body)
  text := string(body)
  defer resp.Body.Close()

  byt := []byte(text)
  var data arraymap
    if err := json.Unmarshal(byt, &data); err != nil {
        panic(err)
    }
    for i, _ := range data {
        tag := data[i]["Key"].(string)
		if tag != fresh_tag {
     resp, err := http.Get("http://172.31.4.66:8500/v1/catalog/service/apache2?tag="+tag)
     if err != nil {
        fmt.Println(err)
    }
      body, _ := ioutil.ReadAll(resp.Body)
	  text := string(body)
      defer resp.Body.Close()

      byt := []byte(text)
      var data arraymap
      var address []string
        if err := json.Unmarshal(byt, &data); err != nil {
            panic(err)
        }
        for i, _ := range data {
            x := data[i]["Address"].(string)
    		address = append(address, x)
        }
    	
    	for item, _ := range address {
    	addr := address[item]
        req, err := http.NewRequest("PUT", "http://"+addr+":8500/v1/agent/service/maintenance/apache2?enable=true&reason=New+deployment+completed", nil)
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
   }
 }

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
  datafinal := string(data)
  response = strings.Trim(datafinal, "\"")
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

//function for creating nodename based on previous entries
func kv_get_nodename(Nodename string, AppVersion string) (string) {
data := kv_get(AppVersion)
var count string
if data != "" {
 r := "-"
 x := strings.LastIndex(data, r) + 1
 Nodenumber, _ := strconv.ParseInt(data[x:], 0, 64)
 Nodenumber++
 Nodenumberstring := strconv.Itoa(int(Nodenumber))
 count = Nodenumberstring
 }else {
   count = "1"
}
 response := Nodename+ "-" + AppVersion + "-" + count
 return response
 }


//main kv function for receiving existing entries and creating new ones
func kv(Nodename string, AppV string) (string) {
 data := kv_get(AppV)
NewNodename := kv_get_nodename(Nodename, AppV)
var NodeValues string
if data != "" {
 NodeValues = data + "," + NewNodename
}else {
 NodeValues = NewNodename
}
kv_put(NodeValues, AppV)
return NewNodename
}


func Loop(w http.ResponseWriter, r *http.Request) {
r.ParseForm()
Nodestring := r.FormValue("nodename")
Imagestring := r.FormValue("image")
Countint, _ := strconv.ParseInt(r.FormValue("inscount"), 0, 64)
AppV := r.FormValue("version")

for Countint > 0 {
   Nodestringfinal := kv(Nodestring, AppV)
   output := CreateInstance(Nodestringfinal, Imagestring, AppV)
   Countint--
   fmt.Fprintf(w, output)
              }
  time.Sleep(80 * time.Second)
  maint(AppV)
}

func CreateInstance(Nodename string, Image string, AppV string) (string) {
data := "#!/bin/bash \n ip=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1) \n sed -i s/enter_node_name_here/"+ Nodename  +"/g /etc/consul.d/agent.json \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/agent.json \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/apache.json \n sed -i s/enter_tag_here/"+ AppV +"/g /etc/consul.d/apache.json \n sed -i s/ENTER_IP_HERE/$ip/g /var/www/html/index.html \n consul agent -config-dir /etc/consul.d/"
sEnc := b64.StdEncoding.EncodeToString([]byte(data))

svc := ec2.New(session.New(&aws.Config{Region: aws.String("us-west-1")}))
    runResult, resp := svc.RunInstancesRequest(&ec2.RunInstancesInput{
        ImageId:      aws.String(Image),
        InstanceType: aws.String("t2.micro"),
        KeyName:      aws.String("go"),
        MinCount:     aws.Int64(1),
        MaxCount:     aws.Int64(1),
        UserData:     aws.String(sEnc),
        SubnetId:     aws.String("subnet-2fc31377"),
        SecurityGroupIds: aws.StringSlice([]string{"sg-8cf887f5"}),
        TagSpecifications: []*ec2.TagSpecification{
         &ec2.TagSpecification{
                ResourceType: aws.String(ec2.ResourceTypeInstance),
                Tags: []*ec2.Tag{
                 &ec2.Tag{
                  Key:   aws.String("Name"),
                  Value: aws.String(Nodename),
                 },
   },
},
},
},
)

    err := runResult.Send()
	var msg *ec2.Reservation
  if err == nil {
    msg = resp
    }	
  json_val, _ := json.MarshalIndent(msg, "", " ")
 return string(json_val)
}


var htmlStr = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
</head>
<body>
  <div>
      <form  action="/createins">
      Node Name: <input type="text" name="nodename" value="App" > <br>
      Version: <input type="text" name="version" value="" > <br>
      Image:  <input type="text" name="image" value="ami-6eeefd0e" > <br>
      Instance Count:  <input type="number" name="inscount" value="" > <br>
          <input type="submit" value="submit" />
      </form>


  </div>
</body>
</html>
`


