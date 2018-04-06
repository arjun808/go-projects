package main

import (
   "github.com/aws/aws-sdk-go/aws"
   "github.com/aws/aws-sdk-go/aws/session"
   "github.com/aws/aws-sdk-go/service/ec2"
   "net/http"
   "fmt"
   b64 "encoding/base64"
    "strconv"
	"encoding/json"
       )

func main(){
   http.HandleFunc("/",welcome)
   http.HandleFunc("/createins", Loop)
   http.ListenAndServe(":9090",nil)
}

func welcome(w http.ResponseWriter, r *http.Request) {

fmt.Fprintf(w,htmlStr)

}



func Loop(w http.ResponseWriter, r *http.Request) {
r.ParseForm()
Nodestring := r.FormValue("nodename")
Imagestring := r.FormValue("image")
Countint, _ := strconv.ParseInt(r.FormValue("inscount"), 0, 64)
i, _ := strconv.ParseInt("0", 0, 64)
for Countint > i {
   Nodestringfinal := fmt.Sprintf("%s-%d", Nodestring, i)
   output := CreateInstance(Nodestringfinal, Imagestring, Countint, i)
   i++
   fmt.Fprintf(w, output)
              }
}

func CreateInstance(Nodename string, Image string, Inscount int64, i int64) (string) {
data := "#!/bin/bash \n ip=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1) \n echo $ip > /tmp/ip.txt \n sed -i s/App4.3/"+ Nodename  +"/g /etc/consul.d/agent.json \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/agent.json \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/apache.json \n sed -i s/3/$ip/g /var/www/html/index.html \n consul agent -config-dir /etc/consul.d/"
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
      Image:  <input type="text" name="image" value="ami-4f89982f" > <br>
      Instance Count:  <input type="number" name="inscount" value="" > <br>
          <input type="submit" value="submit" />
      </form>


  </div>
</body>
</html>
`


