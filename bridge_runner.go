//go:build runner

package main

import (
	"log"
)

func main() {
	app := NewApp()
	log.Println("Starting Dev RPC Bridge on 127.0.0.1:3999...")
	if err := StartDevBridge(app, "127.0.0.1:3999"); err != nil {
		log.Fatalf("Bridge error: %v", err)
	}
}
