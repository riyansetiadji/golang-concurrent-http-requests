package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 200, // Match with numWorkers
	},
}

func sendRequest(input int) string {
	var resp *http.Response
	var req *http.Request
	var err error

	if input%2 == 0 {
		//POST
		url := "http://localhost:3030/users"
		data := fmt.Sprintf(`{"index": %d}`, input)
		var jsonStr = []byte(data)

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		resp, err = client.Do(req)
		//resp, err = client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	} else {
		//GET
		resp, err = client.Get("http://localhost:3030/")
	}

	if err != nil {
		return err.Error()
	}
	io.Copy(ioutil.Discard, resp.Body)
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

	numCalls := 100000
	numWorkers := 200

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
