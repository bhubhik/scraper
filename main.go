package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var urls = []string{"https://www.imdb.com/chart/top/", "https://www.imdb.com/chart/toptv/", "https://www.imdb.com/search/title/?genres=documentary"}

func main() {
	done := make(chan string)

	for _, url := range urls {
		go func(u string) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			fmt.Println("Fetching for: ", u)
			getBody(ctx, u)
			done <- u + " done"
		}(url)
	}
	for range urls {
		msg := <-done
		fmt.Println("Finished: ", msg)
	}
}

func getBody(ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Trouble making request: ", err)
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return ""
	}
	bodyStr := string(body)
	if len(bodyStr) > 100 {
		bodyStr = bodyStr[:100]
	}

	fmt.Printf("Response is: %d\nBody is: %s\n", resp.StatusCode, bodyStr)

	return string(body)
}
