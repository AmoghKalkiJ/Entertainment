package main

import (
	types "github.com/AmoghKalkiJ/Entertainment/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)


const(
	mongoURL="ds135852.mlab.com:35852"
		database = "cxcv"
username= "govind"
password= "govind123"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Hello")
	r.HandleFunc("/novels", NovelsHandler).Methods("POST")
	r.HandleFunc("/movies", MoviessHandler).Methods("POST")
	//r.HandleFunc("/authorlist", AuthorListHandler).Methods("GET")
	fmt.Println("Started server on localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func NovelsHandler(w http.ResponseWriter, r *http.Request) {
	dialogflowreq := types.DialogFlowRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Got some error in reading request body", err)
	}
	fmt.Println("Got body in novels:",string(body))
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	mongoresponse:=GetNovelFromMongo(r,dialogflowreq)
        fmt.Printf("Got mongo response : %+v\n ",mongoresponse)
	author := "Agatha Cristie"
	if dialogflowreq.QueryResult.Parameters.Author != "" {
		author = dialogflowreq.QueryResult.Parameters.Author
	}

	response := types.DialogFlowResponse{
		FulfillmentText:  "Hi you are searching for " + author,
		Source:      "Entertainment app",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)

}

func GetNovelFromMongo(r *http.Request, userrequest types.DialogFlowRequest)types.MongoDBResponse{
	//var sessionMongo *mgo.Session
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoURL},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Printf("CreateSession: %s\n", err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(database).C("books")
        author:= userrequest.QueryResult.Parameters.Author
	category:= userrequest.QueryResult.Parameters.Category
	var book types.MongoDBResponse
	err = collection.Find(bson.M{"authorname": author, "genre": category}).One(&book)
	if err != nil {
		fmt.Printf("DBquery: %s\n", err)
	}
	fmt.Printf("Got book : %+v\n",book)
  return book
}


func MoviessHandler(w http.ResponseWriter, r *http.Request) {
	dialogflowreq := types.DialogFlowRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Got some error in reading request body", err)
	}
	fmt.Println("Got body in movies:",string(body))
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}
        // Call the movies api
	Url:= "https://api.themoviedb.org/3/discover/movie?api_key=d88ce396d9ce2fe1594d82569fb87955&with_genres=18&primary_release_year=2018"
        movieresponse,err :=http.Get(Url)
	if err != nil {
		fmt.Println("Got some error in movies get api call", err)
	}

	movies:=types.Movies{}
	movieslist , err:= ioutil.ReadAll(movieresponse.Body)

	if err != nil {
		fmt.Println("Got some error in reading movies body ", err)
	}

	err = json.Unmarshal(movieslist, &movies)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	var moviedetails []types.MovieResult
	for i:=0;i<5;i++{
		moviedetails=append(moviedetails,movies.Results[i])
	}

	response := types.DialogFlowResponse{
		FulfillmentText: "Hi please find the below list of movies",
		Source:      "Entertainment app",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)
}


