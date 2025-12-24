package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	url     string
	status  int
	healthy bool
}

func check(url string, results chan CheckResult) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Get(url)

	if err != nil {
		results <- CheckResult{url: url, status: 0, healthy: false}
		return
	}

	defer res.Body.Close()
	
	healthy	:= res.StatusCode == 200

	results <- CheckResult{url: url, status: res.StatusCode, healthy: healthy}
}

func main() {
	fmt.Println("Iniciando monitoramento")

	var wg sync.WaitGroup
	results := make(chan CheckResult)

	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		}

	for _, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()
			check(u, results)
		}(url)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.healthy {
			fmt.Printf(("[UP]   %s (Status: %d)\n"), result.url, result.status)
		} else {
			fmt.Printf(("[DOWN] %s (Status: %d)\n"), result.url, result.status)
		}
	}
	

	fmt.Println("Finalizando monitoramento")
}