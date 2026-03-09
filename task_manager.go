package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

type AuthUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

var db *sql.DB
var jwtSecret = []byte("your-secret-key")

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser := r.Header.Get("Authorization")
		if authUser == "" {
			http.Error(w, "Unauthorized or missing token.\n", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authUser, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func main() {
	var err error
	//this connects my api with the database(matching credential)//
	db, err = sql.Open("postgres", "postgres://postgres:Syarronbeast%2E@localhost:5432/user_api?sslmode=disable")

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

	//function handler for the middlerware

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var authorize AuthUser
			json.NewDecoder(r.Body).Decode(&authorize)
			hashedpassword, err := bcrypt.GenerateFromPassword([]byte(authorize.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error hashing password", http.StatusInternalServerError)
				return
			}
			db.QueryRow("INSERT INTO users_jwt(username , password)VALUES( $1, $2) RETURNING id", authorize.Username, string(hashedpassword)).Scan(&authorize.Id)
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(authorize)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// what user sent
		var clientUser AuthUser
		json.NewDecoder(r.Body).Decode(&clientUser)

		// fetch from database using the username user sent
		var dbUser AuthUser
		db.QueryRow("SELECT id, username, password FROM users_jwt WHERE username=$1", clientUser.Username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password)

		// compare passwords
		err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(clientUser.Password))
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": dbUser.Id,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	})

	//GET METHOD , meaning json encoding and sending to client //
	http.HandleFunc("/tasks", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
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

	}))

	http.HandleFunc("/tasks/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	fmt.Println("Server running on port 9090")
	http.ListenAndServe(":9090", nil)
}
