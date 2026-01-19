package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var host = flag.String("host", "https://api-hookground.sonht.io.vn", "server hostname to use")

	// server args
	var key = flag.String("t", "", "token to identify the user")
	var targetHost = flag.String("h", "", "Target host into which the webhook goes")
	flag.Parse()

	if *key == "" {
		log.Fatal("Token is required. Use -t flag to input your token.")
	}

	if *targetHost == "" {
		log.Fatal(
			"Target host is required. Use -h flag to input your target host at you local.\nFor example: -h localhost:8080/path",
		)
	}

	// client usage: groktunnel [-h=<server hostname>] <local port>
	poll(*host, *key, *targetHost)
}

func poll(host, key, targetHost string) {
	fmt.Printf("💪🏼 Ready to receive webhook from server\n")
	fmt.Printf("🧐 Waiting for webhooks. If there are no webhooks for 5 minutes, the program will exit.\n\n")

	const timeout = 5 * time.Minute
	lastPayloadTime := time.Now()

	for {
		// Check if timeout has been reached
		if time.Since(lastPayloadTime) > timeout {
			fmt.Printf("\n⏱️  No webhooks received for 5 minutes. Shutting down...\n")
			os.Exit(0)
		}

		resp, err := http.Get(host + "/poll/" + key)
		if err != nil {
			log.Println("Error polling server:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println("Error polling server, status code:", resp.StatusCode)
		}

		if len(body) == 0 || resp.StatusCode != 200 {
			time.Sleep(2 * time.Second) // No request, retry after 2s
			continue
		}

		// Parse the response to extract payload if it exists
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			if payload, ok := data["payload"]; ok {
				// Forward only the payload field
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					log.Println("Error marshaling payload:", err)
					continue
				}
				body = payloadBytes
			}
		}

		fmt.Printf("🚚 Received webhook from server: \n%v\n\n", string(body))
		fmt.Printf("🛫 Forwarding to %v\n\n", targetHost)
		forward(body, &targetHost)

		// Update last payload time after successful forward
		lastPayloadTime = time.Now()
	}
}

func forward(payload []byte, targetHost *string) {
	resp, err := http.Post(*targetHost, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("🚨 Error forwarding webhook:", err)
		fmt.Println(
			"Make sure the target host is correct and the server is running, also it should support POST method",
		)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("🚨 Error forwarding webhook:", resp.Status)
		return
	}

	fmt.Printf("✅ Done\n")
	fmt.Printf("--------------------------------\n\n")
}
