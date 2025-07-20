package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bhubhik/scraper/scraper"
)

var urls = []string{"https://www.imdb.com/chart/top/", "https://www.imdb.com/chart/toptv/", "https://www.imdb.com/search/title/?genres=documentary"}

func main() {
	done := make(chan []string)

	for _, url := range urls {
		go func(u string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			fmt.Println("Fetching for: ", u)
			titles := getBody(ctx, u)
			done <- titles
		}(url)
	}

	for range urls {
		titles := <-done
		fmt.Println("Titles:", titles)
	}
}

func getBody(ctx context.Context, url string) []string {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Trouble making request: ", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := scraper.ParseTop10(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil
	}

	return body
}
