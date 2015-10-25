package main

import (  
    // Standard library packages
	"encoding/json"
    "fmt"
    "net/http"
    // Third party packages
    "github.com/julienschmidt/httprouter"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "strconv"
    "net/url"
    "math/rand"
    "time"
)
type Request struct {
     // Id bson.ObjectId `json:"id" bson:"_id"`
	    
       Name   string `json:"name" bson:"name"`
       Address string `json:"address" bson:"address"`
       City string `json:"city" bson:"city"`   
       State string `json:"state" bson:"state"`
       Zip   string    `json:"zip" bson:"zip"` 
       
       
}
var mgoSession *mgo.Session 

type Response struct
{
      // Id bson.ObjectId `json:"_id "bson:_id"`
       Id int `json:"id" "bson":"id"`
       Name   string `json:"name" bson:"name"`
       Address string `json:"address" bson:"address"`
       City string `json:"city" bson:"city"`   
       State string `json:"state" bson:"state"`
       Zip    string    `json:"zip" bson:"zip"`
       Coordinates interface{} `json:"coordinates" bson:"cooridnates"`
      

    
}
/*
type UserController struct{}

func NewUserController() *UserController {  
    return &UserController{}
}
*/

/*type Controller struct {  
    session *mgo.Session
}*/

/*func getSession() *mgo.Session {  
    // Connect to our local mongo
    return s
}*/


func postt(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
   // fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("name"))
    
    //session = getSession()
    u:= Request{}
    resp := Response{}
    //lat float64
    json.NewDecoder(r.Body).Decode(&u)
    address:=u.Address + " " + u.City + " " + u.State + " " + u.Zip
    //fmt.Println(address)
    lat:=returnlatlng(address)

    //latitude, _ := lat.(float64);

    //longitude, _ := lng.(float64);
    
    //fmt.Println("The latittude",lat,"and the longitude",lng)
    mgoSession, err := mgo.Dial("mongodb://goassignment:goassignment@ds043714.mongolab.com:43714/goassignment")

    // Check if connection error, is mongo running? 
    if err != nil {
        panic(err)
    }
  
    // Populate the user data
    rand.Seed(time.Now().UTC().UnixNano())
    counter := rand.Intn(5000)
    resp.Id = counter
    
    resp.Name = u.Name
    resp.Address = u.Address
    resp.City = u.City
    resp.State = u.State
    resp.Zip = u.Zip
    resp.Coordinates = lat
    //resp.Longitudes = longitude
     //uc.session.DB("go_rest_tutorial").C("users").Insert(u)
    // Marshal provided interface into JSON structure
    mgoSession.DB("goassignment").C("users").Insert(resp)
    uj, _ := json.Marshal(resp)
    /*if err := mgoSession.DB("go_rest_tutorial").C("users").FindId(oid).One(&u); err != nil {
        w.WriteHeader(404)
        return
    }*/ 
    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
    fmt.Fprintf(w, "%s", uj)
    
    
}

func returnlatlng(address string) (interface{}){
    var Url *url.URL
    Url, err := url.Parse("http://maps.google.com")
    if err != nil {
      panic("Error Panic")
    }

    Url.Path += "/maps/api/geocode/json"
    //fmt.Println(address)
    parameters := url.Values{}
    parameters.Add("address", address)
    Url.RawQuery = parameters.Encode()
    Url.RawQuery += "&sensor=false"

   // fmt.Println("URL " + Url.String())

    res, err := http.Get(Url.String())
                  //
    if err != nil {
      panic("Error Panic")
    }
  //  fmt.Println(res)
    defer res.Body.Close()
    var v map[string] interface{}
    dec:= json.NewDecoder(res.Body);
    if err := dec.Decode(&v); err != nil {
      fmt.Println("ERROR: " + err.Error())
    }   

    lat := v["results"].([]interface{})[0].(map[string] interface{})["geometry"].(map[string] interface{})["location"]//.(map[string]interface{})["lat"]
    //fmt.Println(lat)
   // lng := v["results"].([]interface{})[0].(map[string] interface{})["geometry"].(map[string] interface{})["location"].(map[string]interface{})["lng"]
//fmt.Println("The latittude",lat,"and the longitude",lng)
    return lat
}


func gett(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
   // fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("name"))
      id := p.ByName("id")
      number, _:= strconv.Atoi(id)
      
    mgoSession, err := mgo.Dial("mongodb://goassignment:goassignment@ds043714.mongolab.com:43714/goassignment")
// https://api.mongolab.com/api/1/databases?apiKey=NehIyTKy-1dStg0RKySzPjAWpKd39ful
    //mongodb://goassignment:goassignment@ds043714.mongolab.com:43714/goassignment
    // Check if connection error, is mongo running? 
    if err != nil {
        panic(err)
    }
    // Stub user
    u := Response{}

    // Fetch user
    if err := mgoSession.DB("goassignment").C("users").Find(bson.M{"id":number}).One(&u); err != nil {
        w.WriteHeader(404)
        return
    }
   /* query := mgoSession.DB("go_rest_tutorial").C("users").Find(bson.M{"id":p.ByName("id")})
    fmt.Println(query)*/
    // Marshal provided interface into JSON structure
    uj, _ := json.Marshal(u)

    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}

func delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {  
    // TODO: only write status for now
 id := p.ByName("id")
 number, _:= strconv.Atoi(id)
    

      
    mgoSession, err := mgo.Dial("mongodb://goassignment:goassignment@ds043714.mongolab.com:43714/goassignment")
    // Check if connection error, is mongo running? 
    if err != nil {
        panic(err)
    }
    // Stub user
   // u := Response{}

    if err := mgoSession.DB("goassignment").C("users").Remove(bson.M{"id":number}); err != nil {
        w.WriteHeader(404)
        return
    }

    // Write status
    w.WriteHeader(200)
}



func put(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  
  id := p.ByName("id")
  number, _:=strconv.Atoi(id)

    mgoSession, err := mgo.Dial("mongodb://goassignment:goassignment@ds043714.mongolab.com:43714/goassignment");
    if err!=nil{
        panic(err)
    }

    res := Response{}
    req := Request{}

    json.NewDecoder(r.Body).Decode(&req)

    if err := mgoSession.DB("goassignment").C("users").Find(bson.M{"id":number}).One(&res); err!=nil{
        w.WriteHeader(404)
        return
    }
    
    res.Id = number    
    if req.Name != ""{
        req.Name = res.Name
    }
    if req.Address != ""{
        res.Address = req.Address
    }
    if req.City != ""{
        res.City = req.City
    }
    if req.State != ""{
        res.State = req.State
    }
    if req.Zip != ""{
        res.Zip = req.Zip
    }

    address :=res.Address+" "+res.City+" "+res.State+" "+res.Zip
    res.Coordinates = returnlatlng(address)

    if err := mgoSession.DB("goassignment").C("users").Update(bson.M{"id":number}, res); err != nil {
        w.WriteHeader(404)
        return
    }

    data, _ := json.Marshal(res)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(201)
    fmt.Fprintf(w, "%s", data)


}


func main() {  
    // Instantiate a new router
 // uc := controllers.NewUserController(getSession())  

    r := httprouter.New()

//r.GET("/user/:id", gett)
r.POST("/locations", postt)
r.GET("/locations/:id", gett)
r.DELETE("/locations/:id", delete)
r.PUT("/locations/:id", put)
	 server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: r,
    }
    server.ListenAndServe()

 //   http.ListenAndServe("localhost:8080", r)
}