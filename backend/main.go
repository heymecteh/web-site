package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("sqlite", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT OR IGNORE INTO users(username, password) VALUES(?, ?)",
		"testuser", "testpass")
	if err != nil {
		log.Fatal(err)
	}

	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filename))
	frontendDir := filepath.Join(rootDir, "frontend")

	http.Handle("/", http.FileServer(http.Dir(frontendDir)))
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedPass string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPass)

	if err != nil || password != storedPass {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Login successful!")
}
