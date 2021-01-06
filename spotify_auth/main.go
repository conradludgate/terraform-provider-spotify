package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	clientID = "956aed6fce0c49ebb0eb1d050d9223ed"
	scopes   = []string{
		"user-read-email",
		"user-read-private",
		// "playlist-read-collaborative",
		"playlist-read-private",
		"playlist-modify-private",
		"playlist-modify-public",
		"user-library-read",
		"user-library-modify",
		// "ugc-image-upload",
	}
	redirectURI = "http://localhost:27228/spotify_callback"
)

func main() {
	codeVerifier := randString("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.-~", 128)
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.URLEncoding.EncodeToString(hash[0:32])

	state := uuid.New().String()

	authURL, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		log.Fatalf("could not create auth url: %s\n", err)
	}
	query := url.Values{}
	query.Set("client_id", clientID)
	query.Set("response_type", "code")
	query.Set("redirect_uri", redirectURI)
	query.Set("code_challenge", codeChallenge)
	query.Set("code_challenge_method", "S256")
	query.Set("state", state)
	query.Set("scope", strings.Join(scopes, " "))
	authURL.RawQuery = query.Encode()

	if BrowserOpen(authURL.String()) {
		handler := http.NewServeMux()
		server := http.Server{
			Addr:    ":27228",
			Handler: handler,
		}

		handler.HandleFunc("/spotify_callback", func(w http.ResponseWriter, r *http.Request) {
			if r.FormValue("error") != "" {
				log.Fatalf("There was an error from spotify: %s\n", r.FormValue("error"))
			}

			if r.FormValue("state") != state {
				log.Fatalf("state value was not valid: %s\n", err)
			}

			authCode := r.FormValue("code")

			fmt.Printf("successfully retrieved authCode:\n\n%s\n\ncode verifier: %s\n", authCode, codeVerifier)

			os.Exit(0)
		})

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Open this page in your browser to authenticate with spotify.")
		fmt.Println()
		fmt.Println(authURL)
		fmt.Println()
		fmt.Println("Once you authenticate, copy the url and paste it here:")
		fmt.Println()

		var returnURLString string
		fmt.Scanln(&returnURLString)
		returnURL, err := url.Parse(returnURLString)
		if err != nil {
			log.Fatalf("provided url was not valid: %s\n", err)
		}

		if returnURL.Query().Get("error") != "" {
			log.Fatalf("There was an error from spotify: %s\n", returnURL.Query().Get("error"))
		}

		if returnURL.Query().Get("state") != state {
			log.Fatalf("state value was not valid: %s\n", err)
		}

		authCode := returnURL.Query().Get("code")

		fmt.Println()
		fmt.Println()
		fmt.Printf("successfully retrieved authCode:\n\n%s\n\ncode verifier: %s\n", authCode, codeVerifier)
	}
}

func randString(charSet string, length int) string {
	max := big.NewInt(int64(len(charSet)))
	output := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			log.Fatalf("Error creating random string: %s\n", err)
		}

		output = append(output, charSet[int(n.Int64())])
	}

	return string(output)
}
