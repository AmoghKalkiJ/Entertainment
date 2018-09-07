package main

import (
	"encoding/json"
	"fmt"
	types "github.com/AmoghKalkiJ/Entertainment/types"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	mongoURL = "ds135852.mlab.com:35852"
	database = "cxcv"
	username = "govind"
	password = "govind123"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Hello")
	r.HandleFunc("/novels", NovelsHandler).Methods("POST")
	r.HandleFunc("/movies", MoviessHandler).Methods("GET")
	r.HandleFunc("/authorlist", AuthorListHandler).Methods("POST")
	r.HandleFunc("/score", ScoreHandler).Methods("GET")
	fmt.Println("Started server on localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func NovelsHandler(w http.ResponseWriter, r *http.Request) {
	dialogflowreq := types.LibraryRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Got some error in reading request body", err)
	}
	fmt.Println("Got body in novels:", string(body))
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	mongoresponse := GetNovelFromMongo(dialogflowreq)
	fmt.Printf("Got mongo response : %+v\n ", mongoresponse)
	//author := "Agatha Cristie"
	/*if dialogflowreq.QueryResult.Parameters.Author != "" {
		author = dialogflowreq.QueryResult.Parameters.Author
	}*/
	//textresponse:=createtextresponse(mongoresponse.Title,mongoresponse.Authorname,mongoresponse.Year,mongoresponse.URL)
	//fmt.Println("text response:",textresponse)
	data, err := json.Marshal(mongoresponse)
	if err != nil {
		fmt.Println("Got some error in data marshalling body", err)
	}
	response := types.DialogFlowResponse{
		FulfillmentText: string(data),
		Source:          "Entertainment app",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)

}

func createtextresponse(title, author, year, url string) string {

	// textresponse:="Title: "+title+"\\n"+"Author: "+author+"\\n"+"Year: "+year+"\\n"+"Link: "+url
	textresponse := "Title: " + title + "\\" + "\nAuthor: " + author + "\\" + "\nYear: " + year + "\\" + "\nLink: " + url

	return textresponse
}

func AuthorListHandler(w http.ResponseWriter, r *http.Request) {
	dialogflowreq := types.DialogFlowRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Got some error in reading request body", err)
	}
	fmt.Println("Got body in authorslist:", string(body))
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	mongoresponse := GetAuthorsListFromMongo(dialogflowreq)
	fmt.Printf("Got mongo response : %+v\n ", mongoresponse)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(mongoresponse)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)
}

func GetAuthorsListFromMongo(userrequest types.DialogFlowRequest) []string {
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
	category := userrequest.QueryResult.Parameters.Category
	var mongodatalist []types.MongoDBResponse
	var authorslist []string
	err = collection.Find(bson.M{"genre": category}).Sort("-start").All(&mongodatalist)
	if err != nil {
		fmt.Printf("DBquery: %s\n", err)
	}
	fmt.Printf("Got book : %+v\n", mongodatalist)

	for _, book := range mongodatalist {
		author := book.Authorname
		authorslist = append(authorslist, author)
	}
	return authorslist
}

func GetNovelFromMongo(userrequest types.LibraryRequest) types.MongoDBResponse {

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
	author := userrequest.Entities["author"]
	category := userrequest.Entities["category"]
	var book types.MongoDBResponse
	err = collection.Find(bson.M{"authorname": strings.ToUpper(author), "genre": strings.ToUpper(category)}).One(&book)
	if err != nil {
		fmt.Printf("DBquery: %s\n", err)
	}
	fmt.Printf("Got book : %+v\n", book)
	return book
}

func MoviessHandler(w http.ResponseWriter, r *http.Request) {

	// Call the movies api
	Url := "https://api.themoviedb.org/3/discover/movie?api_key=d88ce396d9ce2fe1594d82569fb87955&with_genres=18&primary_release_year=2018"
	movieresponse, err := http.Get(Url)
	if err != nil {
		fmt.Println("Got some error in movies get api call", err)
	}

	movies := types.Movies{}
	movieslist, err := ioutil.ReadAll(movieresponse.Body)

	if err != nil {
		fmt.Println("Got some error in reading movies body ", err)
	}

	err = json.Unmarshal(movieslist, &movies)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	var moviedetails []types.MovieResult
	for i := 0; i < 5; i++ {
		moviedetails = append(moviedetails, movies.Results[i])
	}

	moviedata, err := json.Marshal(moviedetails)
	if err != nil {
		fmt.Println("Got some error in data marshalling body", err)
	}

	response := types.DialogFlowResponse{
		FulfillmentText: string(moviedata),
		Source:          "Entertainment app",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)
}

func ScoreHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("request:",r)

	matchesapi := "http://cricapi.com/api/matches?apikey=mNkUw7G9QcXVJTOo7pPo4NHZnPc2"
	matchesresponse, err := http.Get(matchesapi)
	if err != nil {
		fmt.Println("Got some error in movies get api call", err)
	}
	var matchlist types.MatchList
	matchdata, err := ioutil.ReadAll(matchesresponse.Body)

	if err != nil {
		fmt.Println("Got some error in reading movies body ", err)
	}

	err = json.Unmarshal(matchdata, &matchlist)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}
	var uniqueids []int

	for _, match := range matchlist.Matches {
		count := 0
		if (strings.ToUpper(match.Team1) == "INDIA" || strings.ToUpper(match.Team2) == "INDIA") && match.MatchStarted {
			uniqueids = append(uniqueids, match.UniqueID)
			count++
		} else if count < 2 && match.MatchStarted {
			uniqueids = append(uniqueids, match.UniqueID)
			count++
		}

	}
	fmt.Println("Got uniqueids:", uniqueids)
	var scores []types.Score
	for _, id := range uniqueids {

		scoreapi := "http://cricapi.com/api/cricketScore?apikey=mNkUw7G9QcXVJTOo7pPo4NHZnPc2&unique_id=" + strconv.Itoa(id)
		scoreresponse, err := http.Get(scoreapi)
		if err != nil {
			fmt.Println("Got some error in movies get api call", err)
		}
		var score types.Score
		scoredata, err := ioutil.ReadAll(scoreresponse.Body)

		if err != nil {
			fmt.Println("Got some error in reading movies body ", err)
		}

		err = json.Unmarshal(scoredata, &score)
		if err != nil {
			fmt.Println("Got some error in unmarshalling request body", err)
		}
		scores = append(scores, score)
	}
	fmt.Printf("Got scores: %+v\n", scores)
	scoreslist, err := json.Marshal(scores)
	if err != nil {
		fmt.Println("Got some error in data marshalling body", err)
	}

	response := types.DialogFlowResponse{
		FulfillmentText: string(scoreslist),
		Source:          "Entertainment app",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)
}
