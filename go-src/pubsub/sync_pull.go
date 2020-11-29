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

// [START pubsub_subscriber_sync_pull]
import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/pubsub"
)

func pullMsgsSync(wg *sync.WaitGroup, projectID, subID string) error {
	fmt.Printf("Got message")
	defer wg.Done()
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
	// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
	// the max number of messages the client will hold in memory at a time.
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1000

	// Receive messages for 5 seconds.
	ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			fmt.Printf("Got message :%q\n", string(msg.Data))
			msg.Ack()
		}
	}()

	// Receive blocks until the passed in context is done.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil && status.Code(err) != codes.Canceled {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}

func main() {
	var wg sync.WaitGroup
	//ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "iwasnothing-self-learning"

	// Sets the id for the new topic.
	topicID := "pullperf"
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go pullMsgsSync(&wg, projectID, topicID)
	}
	wg.Wait()
}

// [END pubsub_subscriber_sync_pull]
