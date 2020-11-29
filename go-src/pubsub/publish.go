// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// [START pubsub_quickstart_publisher]
import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

type Message1 struct {
	Name    string `json:"name"`
	Numbers []int  `json:"numbers"`
}

func publish(wg *sync.WaitGroup, projectID, topicID, msg string) error {
	defer wg.Done()
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	//fmt.Printf("Published a message")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	for i := 0; i < 5000; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte(msg),
		})
		// Block until the result is returned and a server-generated
		// ID is returned for the published message.
		id, err := result.Get(ctx)
		if err != nil {
			return fmt.Errorf("Get: %v", err)
		}
		fmt.Printf("Published a message; msg ID: %v\n", id)
	}
	return nil
}

func main() {
	//ctx := context.Background()
	var wg sync.WaitGroup
	var numbers []int
	n := 10000
	// Sets your Google Cloud Platform project ID.
	projectID := "iwasnothing-self-learning"

	// Sets the id for the new topic.
	topicID := "goperf"

	for i := 0; i < n; i++ {
		numbers = append(numbers, i)
	}
	msg1 := &Message1{
		Name:    "matthew",
		Numbers: numbers}
	msg1B, _ := json.Marshal(msg1)
	fmt.Println(string(msg1B))

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go publish(&wg, projectID, topicID, string(msg1B))
	}
	wg.Wait()
}

// [END pubsub_quickstart_publisher]
