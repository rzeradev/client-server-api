package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rzeradev/client-server-api/server/repository"
)

const (
	apiURL       = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	serverPort   = ":8080"
	timeoutFetch = 200 * time.Millisecond
	timeoutDB    = 10 * time.Millisecond
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func fetchCotacao(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch cotacao: %s", resp.Status)
	}

	var cotacao Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return "", err
	}

	return cotacao.USDBRL.Bid, nil
}

func cotacaoHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeoutFetch)
		defer cancel()

		bid, err := fetchCotacao(ctx)
		if err != nil {
			log.Printf("error fetching cotacao: %v", err)
			if ctx.Err() == context.DeadlineExceeded {
				logToFile("Error: Context timeout of 200ms exceeded to obtain the Cotacao")
				println("Error: Context timeout of 200ms exceeded to obtain the Cotacao", err)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctxDB, cancelDB := context.WithTimeout(r.Context(), timeoutDB)
		defer cancelDB()

		if err := repo.SaveCotacao(ctxDB, bid); err != nil {
			log.Printf("error saving cotacao: %v", err)
			logToFile("error saving cotacao: " + err.Error())
			http.Error(w, "error saving cotacao", http.StatusInternalServerError)
			return
		}

		// Uncomment the line below so the client throws 300ms Timeout Context Error
		// time.Sleep(300 * time.Millisecond)

		response := map[string]string{"bid": bid}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func logToFile(message string) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(message)
}

func main() {
	repo, err := repository.InitDB()
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	http.HandleFunc("/cotacao", cotacaoHandler(repo))
	log.Printf("Server running on port %s", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
