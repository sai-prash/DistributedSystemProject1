package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "io/ioutil"
    "encoding/json"
    "bytes"
    // "context"

    // "go.mongodb.org/mongo-driver/mongo"
    // "go.mongodb.org/mongo-driver/mongo/options"
)

type Rooms struct{
	Rooms []RoomType `json: "Rooms"`
}

type RoomType struct{
	Type string `json: "type"`
	Dates []string `json: "Dates"`
}

type Booking struct{
	RoomType string `json: "RoomType"`
	Dates []string `json: "Dates"`
}


func main() {

   http.HandleFunc("/", rootHandler)
   http.HandleFunc("/index.html", serveFile2)
   http.HandleFunc("/getReservation", getReservation)
   http.HandleFunc("/makeReservation", checkandUpdateReservation)
   http.HandleFunc("/updateNodes", updateNodes)
   fmt.Println("Server started at port 8080")
   log.Fatal(http.ListenAndServe(":8080", nil))
}

// var collection *mongo.Collection
// var ctx = context.TODO()


// func connection() {
//     clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
//     client, err := mongo.Connect(ctx, clientOptions)
//     if err != nil {
//         log.Fatal(err)
//     }

//     err = client.Ping(ctx, nil)
//     if err != nil {
//         log.Fatal(err)
//     }

//     collection = client.Database("tasker").Collection("tasks")
//     log.Println("isfksbvskbvsdkfv")

//     quickstartDatabase := client.Database("HotelReservation")
//     podcastsCollection := quickstartDatabase.Collection("room")

//     defer client.Disconnect(ctx)

// }

func rootHandler(w http.ResponseWriter, r *http.Request) {

    http.ServeFile(w, r, "./index.html")
}

func getReservation(w http.ResponseWriter, r *http.Request){
    log.Println("dsdsdsdssdds")
    //connection()
	jsonFile, _ := os.Open("./data.json")
    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)
    w.Header().Set("Content-Type", "application/json")
    jsonResp, _ := json.Marshal(result)
    w.Write(jsonResp)
    
}


func checkandUpdateReservation(w http.ResponseWriter, r *http.Request){
    
    fmt.Println("INSIDE CHECK N UPDATE")
    NodeList := [3]string{"http://node3:8080","http://node2:8080", "http://node1:8080"}

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    for i:=0; i<3;i++{
        fmt.Println(NodeList[i] + "/updateNodes")
        req, err := http.NewRequest("POST", NodeList[i] + "/updateNodes", bytes.NewBuffer([]byte(body)))
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }


        fmt.Println("response Status:", resp.Status)
        fmt.Println("response Headers:", resp.Header)
    }
    
}


func updateNodes(w http.ResponseWriter, r *http.Request){
    fmt.Println("------------Updating NODES----------")

	jsonFile, _ := os.Open("data.json")
    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var rooms Rooms
    json.Unmarshal(byteValue, &rooms)

    rmType := rooms.Rooms
    
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))
    var b Booking
    err = json.Unmarshal([]byte(body), &b)
    fmt.Println("---", b);

    for i:= 0; i<len(rmType);i++{
    	fmt.Println("King", rmType[i].Type )
    	if rmType[i].Type != b.RoomType{
    		continue
    	}
    	dates := rmType[i].Dates
    	fmt.Println("")
    	for j:=0; j <len(dates);j++{
    		if contains(b.Dates, dates[j]){
    			 dates = append(dates[:j], dates[j+1:]...)
    			 break
    		}
    	}
    	rmType[i].Dates = dates
    }
    rooms.Rooms = rmType
    w.Header().Set("Content-Type", "application/json")
    jsonResp, _ := json.Marshal(rooms)
    fmt.Println(rooms)
    w.Write(jsonResp)
    os.WriteFile("data.json", []byte(string(jsonResp)), 0644)

}

func serveFile(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./index.html")
}

func serveFile2(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./index.html")
}

//Add it to seprate util file
func contains(s []string, t string) bool {
    for _, a := range s {
        if a == t {
            return true
        }
    }
    return false
}
