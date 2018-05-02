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
   "log"
   "flag"
   "html/template"
   "github.com/gorilla/websocket"	
       )
//mux router for handleFunc.

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", Loop)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host    string
	}{
		r.Host,
	}
	homeTempl.Execute(w, &v)
}

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

 type formdata struct {
     Name string
     Version string
	 Image string
	 Inscount string
 }



//function to place old servers in maintenance mode. currently placing apache service alone on maintenance.
func maint(AppVersion string) {
fresh_tag := AppVersion
resp, err := http.Get("http://172.31.1.198:8500/v1/kv/?recurse=true")
if err != nil {
   fmt.Println(err)
}
defer resp.Body.Close()
body, _ := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  byt := []byte(string(body))
  var data []map[string]interface{}
    if err := json.Unmarshal(byt, &data); err != nil {
        panic(err)
    }
    for i, _ := range data {
        tag := data[i]["Key"].(string)
		if tag != fresh_tag {
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
kv_get_resp, err := http.Get("http://172.31.1.198:8500/v1/kv/"+AppVersion)
if err != nil {
       fmt.Println(err)
}else { //do the maint stuff here
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
  req, req_err := http.NewRequest("PUT", "http://172.31.1.198:8500/v1/kv/"+AppVersion, body)
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


//Function that creates new deployment and removes old deployment
func Loop(w http.ResponseWriter, r *http.Request) {
 ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	
	var m formdata
	
	read_err := ws.ReadJSON(&m)
		if read_err != nil {
			fmt.Println("Error reading json.", read_err)
		}
    
   fmt.Printf("Got message: %#v\n", m)	
   
   
   Nodestring := m.Name
   AppV := m.Version
   Imagestring := m.Image
   Countint, _ := strconv.ParseInt(m.Inscount, 0, 64)
  intro := fmt.Sprintf("Creating %d new servers with Application %s of version %s : \n", Countint, Nodestring, AppV )
  if err := ws.WriteMessage(websocket.TextMessage, []byte(intro)); err != nil {
					return
				}
  for Countint > 0 {
     Nodestringfinal := kv(Nodestring, AppV)
     CreateInstance(Nodestringfinal, Imagestring, AppV)
     Countint--
	 instance_out := fmt.Sprintf("Instance %s has been created \n", Nodestringfinal)
	 if err := ws.WriteMessage(websocket.TextMessage, []byte(instance_out)); err != nil {
					return
				}
              }
	deploy_check := "Checking new deployment for stability... \n"
	if err := ws.WriteMessage(websocket.TextMessage, []byte(deploy_check)); err != nil {
					return
				}
	
    time.Sleep(80 * time.Second)
	deploy_check = "New application deployment has been stable for 2 minutes. Placing older versions of the application offline. \n"
	if err := ws.WriteMessage(websocket.TextMessage, []byte(deploy_check)); err != nil {
					return
				}
  maint(AppV)
    if err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("write close:", err)
				return
			}
}

func CreateInstance(Nodename string, Image string, AppV string) (string) {
data := "#!/bin/bash \n ip=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1) \n sed -i s/enter_node_name_here/"+ Nodename  +"/g /etc/consul.d/agent.json \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/agent.json \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/apache.json \n sed -i s/enter_tag_here/"+ AppV +"/g /etc/consul.d/apache.json \n sed -i s/ENTER_IP_HERE/$ip/g /var/www/html/apache2/index.html \n service httpd start \n consul agent -config-dir /etc/consul.d/"
sEnc := b64.StdEncoding.EncodeToString([]byte(data))

svc := ec2.New(session.New(&aws.Config{Region: aws.String("us-west-1")}))
    runResult, resp := svc.RunInstancesRequest(&ec2.RunInstancesInput{
        ImageId:      aws.String(Image),
        InstanceType: aws.String("t2.micro"),
        KeyName:      aws.String("go"),
        MinCount:     aws.Int64(1),
        MaxCount:     aws.Int64(1),
        UserData:     aws.String(sEnc),
        SubnetId:     aws.String("subnet-64ee773f"),
        SecurityGroupIds: aws.StringSlice([]string{"sg-fff1e686"}),
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


const homeHTML = `<!DOCTYPE html>
<html>
<body>

   Node Name: <input type="text" id="nodename" value="App" > <br>
   Version: <input type="text" id="version" value="1.0" > <br>
   Image: <input type="text" id="image" value="ami-ac2e33cc" > <br>
   Instance Count: <input type="number" id="inscount" value="" > <br>
 <input type="button" value="Submit" onclick="myFunction()">
 <pre id="fileData"></pre>
 <input type="button" id="newdeps" value="New Deployment" onClick="document.location.reload(true)" hidden>
<script>
function myFunction() {
      var data = document.getElementById("fileData");
      var name = document.getElementById("nodename").value;
      var version = document.getElementById("version").value;
      var image = document.getElementById("image").value;
      var inscount = document.getElementById("inscount").value;
	  obj = { Name : name, Version : version, Image : image, Inscount : inscount};
      var conn = new WebSocket("ws://{{.Host}}/ws");
	   conn.onopen = function()
               {
	  conn.send(JSON.stringify(obj));
	  
	      };
	   conn.onmessage = function(evt) {
                    console.log('data updated');
                    data.textContent += evt.data; 
					 return false;
                }
	   conn.onclose = function(evt) {
					 var x = document.getElementById("newdeps");
					 x.style.display = "block";
					
                }
}
        </script>

</body>
</html>
`
