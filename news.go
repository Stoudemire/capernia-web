package main

import (
        "database/sql"
        "fmt"
        "os"
        "sync"
        "time"

        _ "github.com/lib/pq"
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
        g_NewsDb     *sql.DB
        g_NewsMutex  sync.Mutex
)

func InitNews() bool {
        databaseURL := os.Getenv("DATABASE_URL")
        if databaseURL == "" {
                g_LogWarn.Print("DATABASE_URL not set, running without news database")
                return true
        }

        var err error
        g_NewsDb, err = sql.Open("postgres", databaseURL)
        if err != nil {
                g_LogErr.Printf("Failed to connect to database: %v", err)
                return false
        }

        err = g_NewsDb.Ping()
        if err != nil {
                g_LogErr.Printf("Failed to ping database: %v", err)
                return false
        }

        g_Log.Print("Connected to PostgreSQL database")

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
                id SERIAL PRIMARY KEY,
                title VARCHAR(255) NOT NULL,
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
                WHERE id = $1
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

        var id int
        err := g_NewsDb.QueryRow(`
                INSERT INTO news (title, content, created_at) 
                VALUES ($1, $2, CURRENT_TIMESTAMP)
                RETURNING id
        `, title, content).Scan(&id)

        if err != nil {
                g_LogErr.Printf("Failed to create news: %v", err)
                return 0, err
        }

        return id, nil
}

func UpdateNews(id int, title string, content string) error {
        g_NewsMutex.Lock()
        defer g_NewsMutex.Unlock()

        if g_NewsDb == nil {
                return fmt.Errorf("database not initialized")
        }

        result, err := g_NewsDb.Exec(`
                UPDATE news 
                SET title = $1, content = $2 
                WHERE id = $3
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

        result, err := g_NewsDb.Exec(`DELETE FROM news WHERE id = $1`, id)
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
                LIMIT $1 OFFSET $2
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
