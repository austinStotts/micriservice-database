package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var connectionString = "user=Austin dbname=chat password=12345678 host=postgres sslmode=disable"

// server for listening and responding to database queries
func main() {
	fmt.Println("Roger Roger")

	http.HandleFunc("/api/create", handleCreate)
	http.HandleFunc("/api/room", handleGetMany)
	http.HandleFunc("/init/db", handleChatInit)
	http.HandleFunc("/run/custom/query", handleCustom)

	http.HandleFunc("/stats", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("> stats")
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

	http.ListenAndServe(":3000", nil)
}

var requests int
var start = time.Now()

type stats struct {
	Requests int
	Time     string
}

type room struct {
	Password string
	Room     string
	Limit    int
}

type post struct {
	ID       int
	Message  string
	Username string
	Room     string
	Created  string
	Type     string
}

type create struct {
	Password string
	Message  string
	Username string
	Room     string
	Created  string
	Type     string
}

// input row into database
func handleCreate(res http.ResponseWriter, req *http.Request) {
	fmt.Println("create request")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	requests++
	// decode request body
	decoder := json.NewDecoder(req.Body)
	var request create
	er := decoder.Decode(&request)
	if er != nil {
		panic(err)
	}

	fmt.Printf("request -> %v %v %v %v %v %v", request.Password, request.Message, request.Username, request.Room, request.Created, request.Type)

	// handle credentials
	if request.Password != "force" {
		reject, err := json.Marshal("error... access denied :(")
		if err != nil {
			panic(err)
		}
		res.Write(reject)
	} else {
		accepted, err := json.Marshal("message created / stored :)")
		if err != nil {
			panic(err)
		}
		queryString := fmt.Sprintf("INSERT INTO messages (message, username, room, created, type) VALUES ('%v', '%v', '%v', '%v', '%v')",
			request.Message, request.Username, request.Room, request.Created, request.Type)
		result, err := db.Query(queryString)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		res.Write(accepted)
	}
}

// get all from a room
func handleGetMany(res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	requests++
	// decode request body
	decoder := json.NewDecoder(req.Body)
	var request room
	er := decoder.Decode(&request)
	if er != nil {
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
		queryString := fmt.Sprintf("SELECT * FROM messages WHERE room = '%v' LIMIT %v",
			request.Room, request.Limit)
		result, err := db.Query(queryString)
		if err != nil {
			panic(err)
		}

		array := make([]post, 0)
		for result.Next() {
			var x post
			err := result.Scan(&x.ID, &x.Message, &x.Username, &x.Room, &x.Created, &x.Type)
			if err != nil {
				panic(err)
			}
			array = append(array, x)
		}
		response, err := json.Marshal(array)
		if err != nil {
			panic(err)
		}
		res.Write(response)
	}
}

func handleChatInit(res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	requests++
	q := req.URL.Query()
	password := q.Get("password")
	if password != "force" {
		reject, err := json.Marshal("error... access denied :(")
		if err != nil {
			panic(err)
		}
		res.Write(reject)
	} else {
		db.Query(`DROP TABLE IF EXISTS chat`)
		result, err := db.Query(`CREATE TABLE messages (
			id SERIAL PRIMARY KEY,
			message TEXT NOT NULL,
			username VARCHAR(48) NOT NULL,
			room VARCHAR(60) NOT NULL,
			created VARCHAR(60) NOT NULL,
			type VARCHAR(24) NOT NULL
		)`)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		response, err := json.Marshal("chat table created")
		if err != nil {
			panic(err)
		}
		res.Write(response)
	}
}

func handleCustom(res http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	requests++
	q := req.URL.Query()
	password := q.Get("password")
	if password != "force" {
		reject, err := json.Marshal("error... access denied :(")
		if err != nil {
			panic(err)
		}
		res.Write(reject)
	} else {
		queryString := q.Get("query")
		result, err := db.Query(queryString)
		if err != nil {
			panic(err)
		}
		response, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		res.Write(response)
	}
}
