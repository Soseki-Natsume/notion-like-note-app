package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3" // sqldriverは実際に使用しないため、importして初期化のみする。
)

var db *sql.DB

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./app.db"
	}
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS note (
	id TEXT PRIMARY KEY, 
	title TEXT NOT NULL, 
	content TEXT NOT NULL DEFAULT '', 
	created_at DATETIME NOT NULL, 
	updated_at DATETIME NOT NULL
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", listHundler)
	http.HandleFunc("/read", readNoteHandler)
	http.HandleFunc("/create", createNoteHandler)
	http.HandleFunc("/update", updateNoteHandler)
	http.HandleFunc("/delete", deleteNoteHandler)
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	/*tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}*/

}
