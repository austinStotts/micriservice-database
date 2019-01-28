package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// server for listening and responding to database queries
// listen on port 8991
// handle get request
// handle post request
// and handle custom queries
func main() {
	fmt.Println("Roger Roger")

	http.HandleFunc("/create", handleCreate)

	http.HandleFunc("/stats", func(res http.ResponseWriter, req *http.Request) {
		requests++
		s := time.Since(start)
		since := fmt.Sprintf("%v", s)
		status := stats{
			Requests: requests,
			Time:     since,
		}
		response, err := json.Marshal(status)
		if err != nil {
			panic(err)
		}
		res.Write(response)
	})

	http.ListenAndServe(":8991", nil)
}

var requests int
var start = time.Now()

type stats struct {
	Requests int
	Time     string
}

type post struct {
	Password string
	Username string
	Message  string
}

func handleCreate(res http.ResponseWriter, req *http.Request) {
	requests++
	// decode request body
	decoder := json.NewDecoder(req.Body)
	var request post
	err := decoder.Decode(&request)
	if err != nil {
		panic(err)
	}

	// handle credentials
	if request.Password != "0511unlock" {
		reject, err := json.Marshal("error... access denied :(")
		if err != nil {
			panic(err)
		}
		res.Write(reject)
	} else {
		accepted, err := json.Marshal("welcome friend :)")
		if err != nil {
			panic(err)
		}
		res.Write(accepted)
	}
}
