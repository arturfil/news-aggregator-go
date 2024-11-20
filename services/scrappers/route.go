package scrappers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/arturfil/aggregator-script/helpers"
	"github.com/arturfil/aggregator-script/types"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store types.ScrapperStore // interface type
}

func NewHandler(store types.ScrapperStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, "Api working...")
	})

	router.Route("/news", func(router chi.Router) {
		router.Get("/newsapi", h.getArticlesAlt)

	})
}

// 431d455ab9af4360a53850dbfcf20cd8
func (h *Handler) getNewsApiArticles(w http.ResponseWriter, r *http.Request) {

    apiKey := os.Getenv("NEWS_API_KEY")
    url := fmt.Sprintf("https://newsapi.org/v2/everything?q=politics&apiKey=%s", apiKey)
    
	res, err := http.Get(url)
	if err != nil {
		helpers.WriteERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
        errMsg := fmt.Sprintf("Body wasn't parsed properly %s", err)
		helpers.WriteERROR(w, http.StatusInternalServerError, errors.New(errMsg))
		return
	}

	var newsResponse types.NewsAPIResponse
	err = json.NewDecoder(r.Body).Decode(&newsResponse)
	if err != nil {
        errMsg := fmt.Sprintf("newsReponse wasn't retrieved %s", err)
		helpers.WriteERROR(w, http.StatusInternalServerError, errors.New(errMsg))
		return
	}

    if len(newsResponse.Articles) > 0 {
        fmt.Printf("First article %s\n", newsResponse.Articles[0])
    }

    fmt.Printf("BODY %s", body)

	helpers.WriteJSON(w, http.StatusOK, newsResponse)
}

func (h *Handler) getArticlesAlt(w http.ResponseWriter, r *http.Request) {

    apiKey := os.Getenv("NEWS_API_KEY")
	keyword := "politics"
	
	// Construct the URL
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", keyword, apiKey)
	
	// Create a new HTTP client
	client := &http.Client{}
	
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	
	// Add headers if needed
	req.Header.Add("User-Agent", "news-api-client")
	
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	
	// Check if the status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: API returned status code %d. Response: %s\n", resp.StatusCode, string(body))
		return
	}
	
	// Parse the JSON response
	var newsResponse types.NewsAPIResponse
	err = json.Unmarshal(body, &newsResponse)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

    var articles []types.Article
	
	// Print the results
	fmt.Printf("Found %d articles\n", newsResponse.TotalResults)
	for _, article := range newsResponse.Articles {
        articles = append(articles, article)
		// fmt.Printf("\nArticle %d:\n", i+1)
		// fmt.Printf("Title: %s\n", article.Title)
		// fmt.Printf("Description: %s\n", article.Description)
		// fmt.Printf("URL: %s\n", article.URL)
		// fmt.Printf("Published At: %s\n", article.PublishedAt)
	}

    helpers.WriteJSON(w, http.StatusOK, articles)
}
