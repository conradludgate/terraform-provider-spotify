package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const id = "SpotifyAuthProxy"

var (
	source oauth2.TokenSource
	apiKey string
	token  string

	config *oauth2.Config
)

func main() {
	host := os.Getenv("SPOTIFY_CLIENT_BASE_URI")
	if host == "" {
		host = "http://localhost:27228"
	}
	hostUrl, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	apiKey = randString(charSet, 64)
	token = randString(charSet, 64)

	authUrl := *hostUrl
	authUrl.Path = path.Join(authUrl.Path, "authorize")
	authUrl.RawQuery = url.Values{"token": {token}}.Encode()

	fmt.Println("APIKey:", apiKey)
	fmt.Println("Token: ", token)
	fmt.Println("Auth:  ", authUrl.String())

	redirectUrl := *hostUrl
	redirectUrl.Path = path.Join(redirectUrl.Path, "spotify_callback")

	config = &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  redirectUrl.String(),
		Endpoint:     spotify.Endpoint,
		Scopes: []string{
			"user-read-email",
			"user-read-private",
			// "playlist-read-collaborative",
			"playlist-read-private",
			"playlist-modify-private",
			"playlist-modify-public",
			"user-library-read",
			"user-library-modify",
			// "ugc-image-upload",
		},
	}

	http.HandleFunc("/authorize", Authorize)
	http.HandleFunc("/api/token", APIToken)
	http.HandleFunc("/spotify_callback", SpotifyCallback)

	log.Fatal(http.ListenAndServe(":27228", nil))
}

// APIToken is the endpoint for refreshing API tokens
func APIToken(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != id || password != apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "APIToken: invalid authorization")
		return
	}

	if source == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "APIToken: token not available")
		return
	}

	token, err := source.Token()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "APIToken: could not retrieve token")
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		fmt.Fprintf(w, "APIToken: could not encode JSON response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println("Token Retrieved")
}

// Authorize takes a user through the auth flow to get a new access token
func Authorize(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("token")
	if key != token {
		fmt.Fprintln(w, "Authorize: you are not authorized to sign in to this platform")
		w.WriteHeader(http.StatusUnauthorized)
	}

	state := randString(charSet, 32)
	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write([]byte(state))
	state += base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	http.Redirect(w, r, config.AuthCodeURL(state), http.StatusSeeOther)
}

// SpotifyCallback handles the redirect from spotify after authorizing
func SpotifyCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.FormValue("error"); err != "" {
		fmt.Fprintf(w, "Could not complete authorization: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		fmt.Fprintln(w, "Could not complete authorization: no auth code present")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	if len(state) < 32 {
		fmt.Fprintln(w, "Could not complete authorization: invalid state value")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	state1 := []byte(state[:32])
	stateMac, err := base64.RawURLEncoding.DecodeString(state[32:])
	if err != nil {
		fmt.Fprintln(w, "Could not complete authorization: invalid state value")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mac := hmac.New(sha256.New, []byte(apiKey))
	mac.Write(state1)
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(stateMac, expectedMac) {
		fmt.Fprintln(w, "Could not complete authorization: invalid state value")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		fmt.Fprintln(w, "Could not complete authorization: invalid auth code")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	source = config.TokenSource(context.Background(), token)

	fmt.Fprintln(w, "Authorization successful")
	fmt.Println("Authorization successful")
}

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
