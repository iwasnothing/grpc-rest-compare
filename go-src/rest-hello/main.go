// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	msg "restmsg"
)

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var m msg.Message1
		err := decoder.Decode(&m)
		if err != nil {
			panic(err)
		}
		log.Println(m.Name)
		//for i,n := range m.Numbers {
		//log.Println(i,n)
		//}
		fmt.Fprintf(w, "Hello %s!\n", m.Name)
	default:
		fmt.Fprintf(w, "Hello %s!\n", name)
	}
}
