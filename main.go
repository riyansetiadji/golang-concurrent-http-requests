package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 200,
	},
}

func sendRequest(n int) string {
	url := "http://localhost:3030/users"
	data := fmt.Sprintf(`{"index": %d}`, n)
	var jsonStr = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	return resp.Status
}

func worker(jobs <-chan int, results chan<- string) {
	for job := range jobs {
		results <- sendRequest(job)
	}
}

func createJobs(num int, jobs chan<- int) {
	for i := 0; i < num; i++ {
		jobs <- i
	}
	close(jobs)
}

func main() {
	jobs := make(chan int)
	results := make(chan string)

	numCalls := 30000
	numWorkers := 50

	start := time.Now()

	for index := 0; index < numWorkers; index++ {
		go worker(jobs, results)
	}
	go createJobs(numCalls, jobs)
	for j := 0; j < numCalls; j++ {
		fmt.Println(<-results)
	}

	elapsed := time.Since(start)
	fmt.Printf("It took %s\n", elapsed)
}
