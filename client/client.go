package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	serverURL    = "http://localhost:8080/cotacao"
	timeoutFetch = 300 * time.Millisecond
)

type Response struct {
	Bid string `json:"bid"`
}

func fetchCotacao(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		logError("error creating request: ", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logError("Error: Context  timeout of 300ms exceeded by the server", ctx.Err())
		} else {
			logError("Error in the request: ", err)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch cotacao: %s", resp.Status)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Bid, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutFetch)
	defer cancel()

	bid, err := fetchCotacao(ctx)
	if err != nil {
		logError("error fetching cotacao: ", err)
		return
	}

	content := fmt.Sprintf("Dólar: %s", bid)
	if err := os.WriteFile("cotacao.txt", []byte(content), 0644); err != nil {
		logError("error writing to file: ", err)
		return
	}

	fmt.Println("Cotacao saved to cotacao.txt")
}

func logError(message string, err error) {
	logMsg := fmt.Sprintf("%s %v\n", message, err)
	fmt.Println(logMsg)
	file, fileErr := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		fmt.Println("Error opening log file:", fileErr)
		return
	}
	defer file.Close()
	file.WriteString(logMsg)
}
