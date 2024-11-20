# Legislative News Aggregator

Go script for aggregating legislative news from various sources and inserting them into PostgreSQL database.

## Prerequisites

- Go 1.19+
- Running instance of [news-server-nodejs](https://github.com/arturfil/news-server-nodejs)
- PostgreSQL database (via Docker from news-server-nodejs)

## Environment Setup

Create `.env`:

```env
DSN=postgresql://postgres:postgres@localhost:5432/legislative_news?sslmode=disable
NEWS_API_KEY=your_api_key
```

## Installation

```bash
go mod download
```

## Usage

1. Ensure news-server-nodejs Docker containers are running:
```bash
# In news-server-nodejs repository (use link above)
docker-compose up -d
```

2. Run aggregator:
```bash
# bulid
make build

# execute binary 
make run

# build & run
make all
```

## Features

- Base64 article ID generation
- Duplicate article detection
- Batch database inserts
- Error handling and logging

## Database Integration

```go
func saveToDatabase(db *sql.DB, articles []types.Article) error {
    query := `
        INSERT INTO articles (
            encoded_id, title, content, description,
            url, published_date, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        ON CONFLICT (title) DO NOTHING
    `
    // ... implementation
}
```

## Error Handling

- Connection failures
- Duplicate entries
- API response errors
- Database transaction errors

## Building

```bash
make build
```

## Future optimizations
- I chose golang beacause it is a language that has concurrency built-in. This means that it is very easy to control multiple-thread processes. I mention this because in the future, we can add more API sources (currently one) and then run the processes in parallel. This will help to speed up the aggregating pipeline.
- Golang also compiles to a binary so it's easier to deploy to a maching without having to run it with an interpreter. And the binary is more memory efficient compared to the nodejs counterpart.


