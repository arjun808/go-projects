package main

import (
 "fmt"
// "html/template"
 "encoding/json"
 "log"
 "time"
 "gopkg.in/mgo.v2"
 "gopkg.in/mgo.v2/bson"
 "net/http"
 "github.com/gorilla/mux"
 "github.com/dgrijalva/jwt-go"
)

const(
  database = "go"
  collection = "users"
  col = "usertokens"
)

type Claims struct{
   Timestamp         time.Time
   jwt.StandardClaims
}


type User struct{
 Name string  `bson:"name" json:"name"`
 Password string    `bson:"password" json:"password"`
}

type Tok struct{
 CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
 Name string  `bson:"name" json:"name"`
 Token string `bson:"token" json:"token"`
}


var user User

func init(){
 session, err := mgo.Dial("172.31.20.164")
 if err != nil {
                panic(err)
        }
fmt.Println("Connection established, yo!")
defer session.Close()
}


func main(){
        r := mux.NewRouter()
        r.HandleFunc("/",welcome)
        r.HandleFunc("/login", Login)
        r.HandleFunc("/signup", Signup)
        if err := http.ListenAndServe(":3000", r); err != nil {
                log.Fatal(err)
         }
}

func welcome(w http.ResponseWriter, r *http.Request){
fmt.Fprintf(w,htmlStr)

// := template.Must(template.ParseFiles("welcome.html"))
  //    t.Execute(w, nil)
}

func Signup(w http.ResponseWriter, r *http.Request){
     err :=json.NewDecoder(r.Body).Decode(&user)
       if err !=nil{
       fmt.Println("error in passing the value")
       }
     session, err := mgo.Dial("")
       if err != nil {
                panic(err)
        }

     c := session.DB(database).C(collection)
     pipe := c.Pipe([]bson.M{{"$match": bson.M{"name":user.Name}}})
     resp :=[]bson.M{}
     k := pipe.All(&resp)
     if k != nil {
     fmt.Println("qqq")
     }

    fmt.Println(len(resp))
   if len(resp)==0{
     c.Insert(&user)
        return     
 }
      if len(resp)!=0{
      fmt.Println("user exists")
      }
   fmt.Println(resp)
}

func Login (w http.ResponseWriter, r *http.Request){
     r.ParseForm()
 fmt.Println("name: ", r.FormValue("name"), "password: ", r.FormValue("password"))
//   err :=json.NewDecoder(r.Body).Decode(&user)
  //    if err !=nil{
    //  fmt.Println("error in passing the value")
      //}
    u1 := r.FormValue("name")
    p1 := r.FormValue("password")     
     session, err := mgo.Dial("172.31.20.164")
        if err != nil {
           panic(err)
        }
 fmt.Println("connection established")
 c := session.DB(database).C(collection)
 d := session.DB(database).C(col)
      index := mgo.Index{
                Key:        []string{"createdAt","name"},
                Unique:     false,
                DropDups:   false,
                Background: true,
                Sparse:     true,
               ExpireAfter: time.Duration(30) * time.Second,
}

err = d.EnsureIndex(index)
        if err != nil {
                panic(err)
        }

//  fmt.Println(r.Form())

    pipe := c.Pipe([]bson.M{{"$match": bson.M{"name":u1,"password":p1}}})
     resp := []bson.M{}
   k := pipe.All(&resp)
     if k != nil {
     fmt.Println("qqq")
     }
   fmt.Println(len(resp))
   if len(resp)==0{
        fmt.Fprintf(w,"username or password is invalid")
     } 
  if len(resp)!=0{
   fmt.Println("token")
//session.DB(database).C(col).Insert(&user)
//res :=User{}
//err :=c.Find(bson.M{"name": user.Name}).One(&res)
// if err!= nil {
//fmt.Println("err")
//}
  // fmt.Println(res.Token)
 // fmt.Println(len(res.Token))
// if len(res.Token)==0{
// var s = Token()
//   c.Update(bson.M{"name": user.Name},bson.M{"$set":bson.M{"token":s}})

 check := d.Pipe([]bson.M{{"$match": bson.M{"name":u1}}})
     resp := []bson.M{}
   k := check.All(&resp)
     if k != nil {
     fmt.Println("qqq")
     }
   fmt.Println(len(resp))
if len(resp)==0{
 var s = Token()
token :=&Tok{
Name: u1,
Token: s,
CreatedAt: time.Now(),
}
 d.Insert(token)
}
 if len(resp)!=0{
   fmt.Println("token is present ")
res :=Tok{}
err :=d.Find(bson.M{"name": u1}).One(&res)
 if err!= nil {
fmt.Println("err")
}
   fmt.Println(res.Token)
fmt.Fprintf(w,"%+v", res.Token)
}

// d.Insert(bson.M{"name": user.N},bson.M{"token":s})
  }
}


var htmlStr = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
</head>
<body>
  <div>
      <form  action="/login">
      Name: <input type="text" name="name" value="arjun" >
      Password: <input type="password" name="password"  value="niharika">
          <input type="submit" value="submit" />
      </form>


  </div>
</body>
</html>
`



func Token() string{
  mySigningKey := []byte("MiNdTrEe")
    claims := Claims{
    time.Now(),
    jwt.StandardClaims{
        ExpiresAt: 150,
    },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
    ss, err := token.SignedString(mySigningKey)
    fmt.Printf("%v", ss)
    if err!=nil{fmt.Println("%v",err)}
    return ss
}

