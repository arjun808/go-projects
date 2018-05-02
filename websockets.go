package main

import (
	"flag"
	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
//	"os"
	//"strconv"
	"time"
	"github.com/gorilla/websocket"
	//"encoding/json"
	"fmt"
)

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

func serveWs(w http.ResponseWriter, r *http.Request) {
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
   
   
   tag := m.Name
   tag1 := m.Version
   tag2 := m.Image
   tag3 := m.Inscount
	 if err := ws.WriteMessage(websocket.TextMessage, []byte(tag)); err != nil {
					return
				}
	time.Sleep(2 * time.Second)
	 if err := ws.WriteMessage(websocket.TextMessage, []byte(tag1)); err != nil {
					return
				}
	time.Sleep(2 * time.Second)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(tag2)); err != nil {
					return
				}
	time.Sleep(2 * time.Second)
	if err := ws.WriteMessage(websocket.TextMessage, []byte(tag3)); err != nil {
					return
				}
	time.Sleep(2 * time.Second)
	if err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, again! \n")); err != nil {
					fmt.Println(err)
				}
    time.Sleep(2 * time.Second)
	if err := ws.WriteMessage(websocket.TextMessage, []byte("I said hello!")); err != nil {
					fmt.Println(err)
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

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
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
                    data.textContent = 'Connection closed';
                }
}
        </script>

</body>
</html>
`
