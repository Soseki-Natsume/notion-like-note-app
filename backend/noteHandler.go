package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var listTmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>notion-like-note-app</title>
	<script type="module" src="/static/main.js"></script>
	<link rel="stylesheet" href="/static/notion-like-note-app.css">
</head>
<body>
<div class="app">

    <!-- サイドバー -->
    <aside class="sidebar">

        <div class="workspace">
            <div class="workspace-icon">N</div>
            <div class="workspace-name">UserName</div>
        </div>

        <div class="sidebar-section">
            <div class="section-title">Notes</div>

            <button id="create-btn" class="create-btn">
                <span>＋</span>
                新規作成
            </button>

            <div id="note-list" class="note-list">
                {{range .}}

                <div class="note-item">
                    <button
                        data-id="{{.ID}}"
                        class="select-note-btn"
                    >
                        <span class="note-name">
                            {{if .Title}}{{.Title}}{{else}}無題{{end}}
                        </span>
                    </button>

                    <button
                        data-id="{{.ID}}"
                        class="delete-btn"
                        title="削除"
                    >
                        ×
                    </button>
                </div>

                {{- end}}
            </div>
        </div>

    </aside>


    <!-- メインコンテンツ -->
    <main class="editor">

        <input type="hidden" id="note-id" value="">

        <div class="editor-inner">

            <input
                type="text"
                id="note-title"
                placeholder="タイトルを入力..."
            >

            <textarea
                id="note-content"
                placeholder="内容を入力..."
            ></textarea>

        </div>

    </main>

</div>	
</body>
</html>
`))

func listHundler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	rows, err := db.Query("SELECT id, title, content FROM note ORDER BY updated_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil {
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

func readNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var n Note
	err := db.QueryRow("SELECT id, title, content FROM note WHERE id = ?", id).Scan(&n.ID, &n.Title, &n.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	newID := uuid.New().String()
	now := time.Now().Format("2006-01-02 15:04:05")

	sqlStr := "INSERT INTO note(id, title, content, created_at, updated_at) VALUES(?, ?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, newID, "", "", now, now)
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
	now := time.Now().Format("2006-01-02 15:04:05")

	sqlStr := "UPDATE note SET title = ?, content = ?, updated_at = ? WHERE id = ?"
	_, err = db.Exec(sqlStr, note.Title, note.Content, now, note.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
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
