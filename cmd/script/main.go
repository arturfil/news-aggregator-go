package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arturfil/aggregator-script/db"
	"github.com/arturfil/aggregator-script/helpers"
	"github.com/arturfil/aggregator-script/scripts"
	"github.com/arturfil/aggregator-script/types"
	_ "github.com/joho/godotenv/autoload"
)

type Store struct {
	db *sql.DB
}

func main() {
    dsn := os.Getenv("DSN")
    db, err := db.NewDatabase(dsn)
    if err != nil {
        log.Fatal("Cannot connect to database", err)
    }

    fmt.Printf("\n**** Script started ****\n")

    newsApiUrlData, err := scripts.GetNewsApiArticles()
    if err != nil {
        log.Fatal("Error getting news articles:", err)
    }

    err = saveToDatabase(db.Client, newsApiUrlData.Articles)
    if err != nil {
        log.Fatal("Error saving articles to db:", err)
    }

    fmt.Printf("\n**** Script Ended ****\n")
}

// articleExists checks if an article with the given encoded_id exists
func articleExists(db *sql.DB, encodedID string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM articles WHERE encoded_id = $1)`
    
    err := db.QueryRow(query, encodedID).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("error checking article existence: %v", err)
    }
    
    return exists, nil
}

func saveToDatabase(db *sql.DB, articles []types.Article) error {
    query := `
        INSERT INTO articles (
            encoded_id,
            title,
            content,
            description,
            url,
            published_date,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
    `

    stmt, err := db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing statement: %v", err)
    }
    defer stmt.Close()

    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("error beginning transaction: %v", err)
    }

    savedCount := 0
    for i, article := range articles {
        // Generate encoded ID from URL
        encodedID := helpers.ConvertURLToBase64ID(article.URL)
        
        // Check if article already exists
        exists, err := articleExists(db, encodedID)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("error checking article existence: %v", err)
        }

        if exists {
            log.Printf("Skipping article %d '%s' - already exists", i, article.Title)
            continue
        }

        // Parse published date
        publishedAt, err := time.Parse(time.RFC3339, article.PublishedDate)
        if err != nil {
            log.Printf("Warning: Could not parse date for article '%s': %v", article.Title, err)
            publishedAt = time.Now() // Fallback to current time
        }

        // Execute the insert
        _, err = tx.Stmt(stmt).Exec(
            encodedID,
            article.Title,
            article.Content,
            article.Description,
            article.URL,
            publishedAt,
        )
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("error inserting article '%s': %v", article.Title, err)
        }

        savedCount++
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("error committing transaction: %v", err)
    }

    log.Printf("\n\nSuccessfully saved %d new articles (skipped %d existing)", 
        savedCount, len(articles)-savedCount)
    
    return nil
}

