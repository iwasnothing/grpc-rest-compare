/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"

	//"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/credentials"
	//pb "google.golang.org/grpc/examples/helloworld/helloworld"
	pb "grpc-hello/greeter_client/helloworld"
)

const (
	address     = "grpc-hello-s3yk5iivva-uc.a.run.app:443"
	defaultName = "Mattnew"
)

var (
	logger     = log.New(os.Stdout, "", 0)
	serverAddr = flag.String("server", "", "Server address (host:port)")
	serverHost = flag.String("server-host", "", "Host name to which server IP should resolve")
	insecure   = flag.Bool("insecure", false, "Skip SSL validation? [false]")
	skipVerify = flag.Bool("skip-verify", false, "Skip server hostname verification in SSL validation [false]")
	duration   = flag.Uint("duration", 10, "duration (in seconds) to stream the time from the server for")
)

func init() {
	flag.Parse()
	log.SetFlags(log.Flags() ^ log.Ltime ^ log.Ldate)
}

func main() {
	var wg sync.WaitGroup
	// Contact the server and print out its response.
	name := defaultName
	//ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
	var deadlineMs = flag.Int("deadline_ms", 2000*1000, "Default deadline in milliseconds.")
	var numbers []int32
	n := 10000
	for i := 0; i < n; i++ {
		numbers = append(numbers, int32(i))
	}
	//for i,n := range numbers {
	//log.Println(i,n)
	//}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go calling(&wg, deadlineMs, name, numbers)
	}
	wg.Wait()
}
func calling(wg *sync.WaitGroup, deadlineMs *int, name string, numbers []int32) {
	defer wg.Done()
	logger.Printf("calling")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*deadlineMs)*time.Millisecond)
	defer cancel()
	var opts []grpc.DialOption
	if *serverAddr == "" {
		log.Fatal("-server is empty")
	}
	if *serverHost != "" {
		opts = append(opts, grpc.WithAuthority(*serverHost))
	}
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: *skipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logger.Printf("failed to dial server %s: %v", *serverAddr, err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	n := 5000
	for i := 0; i < n; i++ {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name, Numbers: numbers})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}
}
