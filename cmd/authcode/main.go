package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ranganath42/authcli"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Failed to run: %v", err)
		os.Exit(1)
	}
}
func run() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("load .env file: %w", err)
	}

	// Configure the client to communicate with the provider.
	authClient := authcli.New(
		authcli.ProviderGitHub,
		os.Getenv("GITHUB_CLIENT_ID"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
		"http://localhost:3001/callback",
		//authcli.WithScopes([]string{"user"}),
	)

	// Save the state to validate the callback.
	state := fmt.Sprintf("%d", rand.Intn(time.Now().Nanosecond()))

	http.HandleFunc("/", rootHandler(authClient.AuthorizationURL(state)))
	http.HandleFunc("/callback", callbackHandler(authClient, state))
	http.HandleFunc("/welcome", welcomeHandler())

	return http.ListenAndServe(":3001", nil)
}
