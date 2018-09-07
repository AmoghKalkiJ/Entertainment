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
	"strings"
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
	r.HandleFunc("/authorlist", AuthorListHandler).Methods("POST")
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
	fmt.Println("Got body in novels:",string(body))
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	mongoresponse:=GetNovelFromMongo(dialogflowreq)
        fmt.Printf("Got mongo response : %+v\n ",mongoresponse)
	//author := "Agatha Cristie"
	/*if dialogflowreq.QueryResult.Parameters.Author != "" {
		author = dialogflowreq.QueryResult.Parameters.Author
	}*/
	//textresponse:=createtextresponse(mongoresponse.Title,mongoresponse.Authorname,mongoresponse.Year,mongoresponse.URL)
	//fmt.Println("text response:",textresponse)
        data,err:= json.Marshal(mongoresponse)
	if err != nil {
		fmt.Println("Got some error in data marshalling body", err)
	}
	response := types.DialogFlowResponse{
		FulfillmentText: string(data) ,
		Source:      "Entertainment app",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)

}


func createtextresponse(title,author,year,url string)string{

       // textresponse:="Title: "+title+"\\n"+"Author: "+author+"\\n"+"Year: "+year+"\\n"+"Link: "+url
	textresponse:="Title: "+title+"\\"+"\nAuthor: "+author+"\\"+"\nYear: "+year+"\\"+"\nLink: "+url

	return textresponse
}

func AuthorListHandler(w http.ResponseWriter, r *http.Request){
	dialogflowreq := types.DialogFlowRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Got some error in reading request body", err)
	}
	fmt.Println("Got body in authorslist:",string(body))
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	mongoresponse:=GetAuthorsListFromMongo(dialogflowreq)
	fmt.Printf("Got mongo response : %+v\n ",mongoresponse)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(mongoresponse)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)
}

func GetAuthorsListFromMongo( userrequest types.DialogFlowRequest)[]string{
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
	category:= userrequest.QueryResult.Parameters.Category
	var mongodatalist []types.MongoDBResponse
	var authorslist []string
	err = collection.Find(bson.M{ "genre": category}).Sort("-start").All(&mongodatalist)
	if err != nil {
		fmt.Printf("DBquery: %s\n", err)
	}
	fmt.Printf("Got book : %+v\n",mongodatalist)

	for _,book:=range mongodatalist{
                author:=book.Authorname
		authorslist = append(authorslist,author)
	}
	return authorslist
}

func GetNovelFromMongo( userrequest types.LibraryRequest)types.MongoDBResponse{

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
        author:= userrequest.Entities["author"]
	category:= userrequest.Entities["category"]
	var book types.MongoDBResponse
	err = collection.Find(bson.M{"authorname": strings.ToUpper(author), "genre":strings.ToUpper(category)}).One(&book)
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


