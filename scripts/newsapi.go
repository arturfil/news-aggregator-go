package scripts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/arturfil/aggregator-script/types"
)

func GetNewsApiArticles() (*types.NewsAPIResponse, error) {
	apiKey := os.Getenv("NEWS_API_KEY")
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=politics&apiKey=%s", apiKey)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var newsResponse types.NewsAPIResponse
	err = json.Unmarshal(body, &newsResponse)
	if err != nil {
        return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

    return &newsResponse, nil
}
