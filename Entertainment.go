package main

import (
	types "github.com/AmoghKalkiJ/Entertainment/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	fmt.Println("Hello")
	r.HandleFunc("/novels", NovelsHandler).Methods("POST")
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
	fmt.Printf("Got Body: +v\n",body)
	err = json.Unmarshal(body, &dialogflowreq)
	if err != nil {
		fmt.Println("Got some error in unmarshalling request body", err)
	}

	author := "Agatha Cristie"
	if dialogflowreq.QueryResult.Parameters.Author != "" {
		author = dialogflowreq.QueryResult.Parameters.Author
	}

	response := types.DialogFlowResponse{
		Speech:      "Hi you are searching for " + author,
		DisplayText: "Hi you are searching for " + author,
		Source:      "Webhook",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	outResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(outResponse)

}
