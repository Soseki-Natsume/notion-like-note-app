package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:") // memoryにテスト用db作成
	if err != nil {
		t.Fatalf("テスト用DBオープン失敗: %v", err)
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
		t.Fatalf("テーブル作成失敗: %v", err)
	}
}

// 1. 新規作成（CREATE: /create）のテスト
func TestCreateNoteHandler(t *testing.T) {
	setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/create", nil)
	rec := httptest.NewRecorder()

	createNoteHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待ステータス: 200, 実際: %d", rec.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}

	if res["id"] == "" {
		t.Error("IDが発行されていません")
	}
}

// 2. 特定ノートの取得（READ: /read?id=...）のテスト
func TestReadNoteHandler(t *testing.T) {
	setupTestDB(t)
	defer db.Close()

	// テスト用データを1件準備
	_, err := db.Exec("INSERT INTO note (id, title, content, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))", "read-test-id", "取得テストタイトル", "取得テスト本文")
	if err != nil {
		t.Fatalf("データ準備失敗: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/read?id=read-test-id", nil)
	rec := httptest.NewRecorder()

	readNoteHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待ステータス: 200, 実際: %d", rec.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("レスポンス解析失敗: %v", err)
	}

	if res["title"] != "取得テストタイトル" || res["content"] != "取得テスト本文" {
		t.Errorf("取得データが不一致です: %#v", res)
	}
}

// 3. 更新（UPDATE: /update）のテスト
func TestUpdateNoteHandler(t *testing.T) {
	setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO note (id, title, content, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))", "update-test-id", "旧タイトル", "旧本文")
	if err != nil {
		t.Fatalf("データ準備失敗: %v", err)
	}

	payload := []byte(`{"id": "update-test-id", "title": "新タイトル", "content": "新本文"}`)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	updateNoteHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待ステータス: 200, 実際: %d", rec.Code)
	}

	// DBの値が書き換わったか確認
	var title, content string
	err = db.QueryRow("SELECT title, content FROM note WHERE id = ?", "update-test-id").Scan(&title, &content)
	if err != nil {
		t.Fatalf("DB検索失敗: %v", err)
	}

	if title != "新タイトル" || content != "新本文" {
		t.Errorf("更新内容が不一致です: title=%s, content=%s", title, content)
	}
}

// 4. 削除（DELETE: /delete?id=...）のテスト
func TestDeleteNoteHandler(t *testing.T) {
	setupTestDB(t)
	defer db.Close()

	_, err := db.Exec("INSERT INTO note (id, title, content, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))", "delete-test-id", "削除用タイトル", "削除用本文")
	if err != nil {
		t.Fatalf("データ準備失敗: %v", err)
	}

	payload := []byte(`{"id": "delete-test-id"}`)
	req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	deleteNoteHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待ステータス: 200, 実際: %d", rec.Code)
	}

	// 削除後にデータが残っていないことを確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM note WHERE id = ?", "delete-test-id").Scan(&count)
	if err != nil {
		t.Fatalf("DB検索失敗: %v", err)
	}

	if count != 0 {
		t.Errorf("レコードが削除されていません（件数: %d）", count)
	}
}
