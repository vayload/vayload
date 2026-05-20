package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vayload/vayload/cmd/cli/client"
)

const N = 10

func benchHTTP() {
	start := time.Now()

	for range N {
		res, err := http.Get("http://localhost:9000/health")
		if err != nil {
			log.Fatal(err)
		}

		_, _ = io.ReadAll(res.Body)
		res.Body.Close()
	}

	fmt.Println("HTTP TCP:", time.Since(start))
}

func benchLocal() {
	c := client.NewClient()
	ctx := context.Background()

	start := time.Now()

	for range N {
		res, err := c.Get("health").Send(ctx)
		if err != nil {
			log.Fatal(err)
		}

		_, _ = io.ReadAll(res.Body)
		res.Body.Close()
	}

	fmt.Println("Unix socket / Pipe:", time.Since(start))
}

func main() {
	benchHTTP()
	benchLocal()
}
