package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error
	//this connects my api with the database(matching credential)//
	db, err = sql.Open("postgres", "postgres://postgres:Syarronbeast.@localhost:5432/user_api?sslmode=disable")

	if err != nil {
		fmt.Println("\nError connecting to the database: ", err)
		return
	}
	err = db.Ping() // db.ping checks if connection is still there alive//
	if err != nil {
		fmt.Println("\nError connecting to the database : ", err)
		return
	}
	fmt.Println("\nConnected Succesfully")

	defer db.Close() //well it will be executed after main does it things //

	//GET METHOD , meaning json encoding and sending to client //
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			rows, err := db.Query("SELECT id,title,description,status,created_at FROM tasks")
			if err != nil {
				fmt.Println("\nError with the query DB:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()
			var tasks []Task
			for rows.Next() {
				var t Task
				rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreatedAt)
				tasks = append(tasks, t)
			}
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(tasks)

		} else if r.Method == "POST" {
			var newtask Task
			json.NewDecoder(r.Body).Decode(&newtask)
			db.QueryRow("INSERT INTO tasks (title,description,status,created_at) VALUES ($1 , $2 ,$3 , $4) RETURNING id", newtask.Title, newtask.Description, newtask.Status, newtask.CreatedAt).Scan(&newtask.Id)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newtask)
		} else {
			http.Error(w, "error occured:  ", http.StatusMethodNotAllowed)
			return
		}

	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			id := r.URL.Path[len("/tasks/"):]
			var Updatetask Task
			json.NewDecoder(r.Body).Decode(&Updatetask)
			db.Exec("UPDATE tasks SET title=$1, description=$2, status=$3, created_at=$4 WHERE id=$5", Updatetask.Title, Updatetask.Description, Updatetask.Status, Updatetask.CreatedAt, id)

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Updatetask)

		} else if r.Method == "DELETE" {
			id := r.URL.Path[len("/tasks/"):]
			db.Exec("DELETE FROM tasks where id=$1", id)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "task deleted successfully")
		}
	})

	fmt.Println("Server running on port 9090")
	http.ListenAndServe(":9090", nil)
}
