package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	msg "rest-hello/client/restmsg"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var numbers []int
	n := 10000
	for i := 0; i < n; i++ {
		numbers = append(numbers, i)
	}
	msg1 := &msg.Message1{
		Name:    "matthew",
		Numbers: numbers}
	msg1B, _ := json.Marshal(msg1)
	fmt.Println(string(msg1B))

	for i := 0; i < 10; i++ {
		wg.Add(1)
		calling(&wg, msg1B)
	}
	wg.Wait()
}
func calling(wg *sync.WaitGroup, msg1B []byte) {
	defer wg.Done()
	for i := 0; i < 5000; i++ {
		resp, err := http.Post("https://rest-hello-s3yk5iivva-uc.a.run.app", "application/json", bytes.NewBuffer(msg1B))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Response status:", resp.Status)

		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 5; i++ {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

}
