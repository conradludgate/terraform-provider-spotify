package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		// default for terraform-provider-spotify
		clientID = "956aed6fce0c49ebb0eb1d050d9223ed"
	}

	scopes := strings.Split(os.Getenv("SPOTIFY_SCOPES"), ",")
	if len(scopes) == 0 {
		// default for terraform-provider-spotify
		scopes = []string{
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
	}

	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	if redirectURI == "" {
		// default for terraform-provider-spotify
		redirectURI = "http://localhost:27228/spotify_callback"
	}

	accessToken, expiresIn, err := auth(clientID, scopes, redirectURI)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not authenticate.", err.Error())
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Authenticated successfully. Token expires in", expiresIn.String())

	fmt.Printf("export SPOTIFY_ACCESS_TOKEN=%s\n", accessToken)
}

func auth(clientID string, scopes []string, redirectURI string) (string, time.Duration, error) {
	state := uuid.New().String()

	authURL, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		return "", 0, fmt.Errorf("could not create auth url: %w", err)
	}
	query := url.Values{}
	query.Set("client_id", clientID)
	query.Set("response_type", "token")
	query.Set("redirect_uri", redirectURI)
	query.Set("state", state)
	query.Set("scope", strings.Join(scopes, " "))
	authURL.RawQuery = query.Encode()

	fmt.Fprintln(os.Stderr, "Open this page in your browser to authenticate with spotify.")
	fmt.Fprintln(os.Stderr, authURL)
	fmt.Fprintln(os.Stderr, "Once you authenticate, copy the url and paste it here:")

	var returnURLString string
	fmt.Scanln(&returnURLString)
	returnURL, err := url.Parse(returnURLString)
	if err != nil {
		return "", 0, fmt.Errorf("provided return url was not valid: %w", err)
	}
	data, err := url.ParseQuery(returnURL.Fragment)
	if err != nil {
		return "", 0, fmt.Errorf("provided return url was not valid: %w", err)
	}
	if data.Get("state") != state {
		return "", 0, fmt.Errorf("invalid state found in url")
	}

	expiresIn, err := time.ParseDuration(data.Get("expires_in") + "s")
	if err != nil {
		return "", 0, fmt.Errorf("could not parse expires_in value")
	}

	return data.Get("access_token"), expiresIn, nil
}
