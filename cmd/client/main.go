package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	pb "github.com/golrice/e-fis/internal/protocal"
	"google.golang.org/protobuf/proto"
)

func main() {
	var serverAddr string
	var key string
	flag.StringVar(&serverAddr, "server", "http://localhost:9999", "server address")
	flag.StringVar(&key, "key", "", "key to get from cache")
	flag.Parse()

	if key == "" {
		log.Fatal("key is required")
	}

	request := &pb.Request{
		NodeName: "scores",
		Key:      key,
	}

	body, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(serverAddr+"/api", "application/octet-stream", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("failed to send request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("server returned non-OK status: %v", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	var response pb.Response
	if err := proto.Unmarshal(responseBody, &response); err != nil {
		log.Fatalf("failed to unmarshal response: %v", err)
	}

	fmt.Printf("Value for key '%s': %s\n", key, string(response.Value))
}
