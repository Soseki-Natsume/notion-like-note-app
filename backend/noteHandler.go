package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const notitle = "無題"

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var listTmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>notion-like-note-app</title>
	<script type="module" src="/static/main.js"></script>
</head>
<body>
    {{range .}}
    <button name="note-title" form=""></button>
    <button name="create" form="">+新規作成</button>
    <button name="delete" form="">-削除</button>
	{{- end}}
    <input type="text" id="note-title" name="note-title" placeholder="タイトル">
    <textarea id="note-content" name="note-content" placeholder="内容を入力..."></textarea>
</body>
</html>
`))

func listHundler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	rows, err := db.Query("SELECT title, content FROM note ORDER BY updated_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.Title, &n.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = listTmpl.Execute(w, notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newID := uuid.New().String()
	now := time.Now()

	sqlStr := "INSERT INTO note(id, title, content, created_at, updated_at) VALUES(?, ?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, newID, notitle, "", now, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id": "%s", "title": "notitle", "content": ""}`, newID)
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now()

	sqlStr := "UPDATE note SET title = ?, content = ?, updated_at = ? WHERE id = ?"
	_, err = db.Exec(sqlStr, note.Title, note.Content, now, note.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "success"}`)
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStr := "DELETE FROM note WHERE id = ?"
	_, err = db.Exec(sqlStr, note.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `{"status": "success"}`)
}
