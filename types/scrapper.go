package types

import "time"

type ScrapperStore interface {
	SaveArticle() error
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID            string    `json:"id"`
	EncodedID     string    `json:"encoded_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Source        Source    `json:"source"`
	Author        string    `json:"author"`
	State         string    `json:"state"`
	Topic         string    `json:"topic"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	PublishedDate string    `json:"published_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NewsAPIResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}
