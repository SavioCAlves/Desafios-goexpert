package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(body)
}
