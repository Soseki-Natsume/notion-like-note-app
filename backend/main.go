package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // sqldriverは実際に使用しないため、importして初期化のみする。
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "note.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", listHundler)
	http.HandleFunc("/read", readNoteHandler)
	http.HandleFunc("/create", createNoteHandler)
	http.HandleFunc("/update", updateNoteHandler)
	http.HandleFunc("/delete", deleteNoteHandler)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}*/

}
