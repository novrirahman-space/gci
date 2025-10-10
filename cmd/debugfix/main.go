package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// HTTP client dengan timeout agar operasi jaringan tidak menggantung
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// fetchURL menjalankan satu unit kerja HTTP secara aman & terkontrol.

func fetchURL(ctx context.Context, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err:= http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("request build error:", url, err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("fetch error:", url, err)
		return
	}

	defer resp.Body.Close()

	// Drain body agar koneksi ditutup rapi (keep-alive tetap sehat)
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		fmt.Println("read body error:", url, err)
		return
	}

	fmt.Println("OK:", url, resp.Status)
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.golang.com",
	}

	// Context opsional untuk shutdown terkontrol
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, u:= range urls {
		u := u
		go fetchURL(ctx, u, &wg)
	}

	// Barrier: tunggu semua goroutine selesai
	wg.Wait()
	fmt.Println("Done")
}