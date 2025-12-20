package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type (
	TNews struct {
		ID        int
		Title     string
		Content   string
		CreatedAt time.Time
	}
)

var (
	g_NewsDb    *sql.DB
	g_NewsMutex sync.Mutex
)

func InitNews() bool {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "tibia.db"
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			g_LogErr.Printf("Failed to create database directory: %v", err)
			return false
		}
	}

	var err error
	g_NewsDb, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		g_LogErr.Printf("Failed to connect to database: %v", err)
		return false
	}

	err = g_NewsDb.Ping()
	if err != nil {
		g_LogErr.Printf("Failed to ping database: %v", err)
		return false
	}

	g_Log.Print("Connected to SQLite3 database")

	if !CreateNewsTables() {
		return false
	}

	return true
}

func ExitNews() {
	if g_NewsDb != nil {
		g_NewsDb.Close()
	}
}

func CreateNewsTables() bool {
	query := `
	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := g_NewsDb.Exec(query)
	if err != nil {
		g_LogErr.Printf("Failed to create news table: %v", err)
		return false
	}

	g_Log.Print("News table ensured")
	return true
}

func GetAllNews() ([]TNews, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return []TNews{}, fmt.Errorf("database not initialized")
	}

	rows, err := g_NewsDb.Query(`
		SELECT id, title, content, created_at 
		FROM news 
		ORDER BY created_at DESC
	`)
	if err != nil {
		g_LogErr.Printf("Failed to query news: %v", err)
		return nil, err
	}
	defer rows.Close()

	var news []TNews
	for rows.Next() {
		var n TNews
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt)
		if err != nil {
			g_LogErr.Printf("Failed to scan news row: %v", err)
			continue
		}
		news = append(news, n)
	}

	return news, nil
}

func GetNewsByDateRange(fromDate time.Time, toDate time.Time) ([]TNews, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return []TNews{}, fmt.Errorf("database not initialized")
	}

	rows, err := g_NewsDb.Query(`
		SELECT id, title, content, created_at 
		FROM news 
		WHERE created_at >= ? AND created_at <= ?
		ORDER BY created_at DESC
	`, fromDate, toDate.AddDate(0, 0, 1))
	if err != nil {
		g_LogErr.Printf("Failed to query news by date range: %v", err)
		return nil, err
	}
	defer rows.Close()

	var news []TNews
	for rows.Next() {
		var n TNews
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt)
		if err != nil {
			g_LogErr.Printf("Failed to scan news row: %v", err)
			continue
		}
		news = append(news, n)
	}

	return news, nil
}

func GetNewsByDateRangePaginated(fromDate time.Time, toDate time.Time, page int, itemsPerPage int) ([]TNews, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return []TNews{}, fmt.Errorf("database not initialized")
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * itemsPerPage

	rows, err := g_NewsDb.Query(`
		SELECT id, title, content, created_at 
		FROM news 
		WHERE created_at >= ? AND created_at <= ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, fromDate, toDate.AddDate(0, 0, 1), itemsPerPage, offset)
	if err != nil {
		g_LogErr.Printf("Failed to query news by date range paginated: %v", err)
		return nil, err
	}
	defer rows.Close()

	var news []TNews
	for rows.Next() {
		var n TNews
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt)
		if err != nil {
			g_LogErr.Printf("Failed to scan news row: %v", err)
			continue
		}
		news = append(news, n)
	}

	return news, nil
}

func GetNewsByDateRangeCount(fromDate time.Time, toDate time.Time) (int, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	var count int
	err := g_NewsDb.QueryRow(`
		SELECT COUNT(*) FROM news 
		WHERE created_at >= ? AND created_at <= ?
	`, fromDate, toDate.AddDate(0, 0, 1)).Scan(&count)
	if err != nil {
		g_LogErr.Printf("Failed to count news by date range: %v", err)
		return 0, err
	}

	return count, nil
}

func GetNewsById(id int) (*TNews, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var n TNews
	err := g_NewsDb.QueryRow(`
		SELECT id, title, content, created_at 
		FROM news 
		WHERE id = ?
	`, id).Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		g_LogErr.Printf("Failed to query news by id: %v", err)
		return nil, err
	}

	return &n, nil
}

func CreateNews(title string, content string) (int, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	result, err := g_NewsDb.Exec(`
		INSERT INTO news (title, content, created_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`, title, content)

	if err != nil {
		g_LogErr.Printf("Failed to create news: %v", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		g_LogErr.Printf("Failed to get last insert id: %v", err)
		return 0, err
	}

	return int(id), nil
}

func UpdateNews(id int, title string, content string) error {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return fmt.Errorf("database not initialized")
	}

	result, err := g_NewsDb.Exec(`
		UPDATE news 
		SET title = ?, content = ? 
		WHERE id = ?
	`, title, content, id)

	if err != nil {
		g_LogErr.Printf("Failed to update news: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		g_LogErr.Printf("Failed to get rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("news not found")
	}

	return nil
}

func DeleteNews(id int) error {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return fmt.Errorf("database not initialized")
	}

	result, err := g_NewsDb.Exec(`DELETE FROM news WHERE id = ?`, id)
	if err != nil {
		g_LogErr.Printf("Failed to delete news: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		g_LogErr.Printf("Failed to get rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("news not found")
	}

	return nil
}

func GetTotalNewsCount() (int, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	var count int
	err := g_NewsDb.QueryRow(`SELECT COUNT(*) FROM news`).Scan(&count)
	if err != nil {
		g_LogErr.Printf("Failed to count news: %v", err)
		return 0, err
	}

	return count, nil
}

func GetNewsPaginated(page int, itemsPerPage int) ([]TNews, error) {
	g_NewsMutex.Lock()
	defer g_NewsMutex.Unlock()

	if g_NewsDb == nil {
		return []TNews{}, fmt.Errorf("database not initialized")
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * itemsPerPage

	rows, err := g_NewsDb.Query(`
		SELECT id, title, content, created_at 
		FROM news 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, itemsPerPage, offset)
	if err != nil {
		g_LogErr.Printf("Failed to query news: %v", err)
		return nil, err
	}
	defer rows.Close()

	var news []TNews
	for rows.Next() {
		var n TNews
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt)
		if err != nil {
			g_LogErr.Printf("Failed to scan news row: %v", err)
			continue
		}
		news = append(news, n)
	}

	return news, nil
}
