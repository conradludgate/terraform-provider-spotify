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
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/caarlos0/env/v6"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const id = "SpotifyAuthProxy"

var (
	source oauth2.TokenSource
	envCfg *EnvConfig
	token  string

	config *oauth2.Config
)

type EnvConfig struct {
	ClientID     string   `env:"SPOTIFY_CLIENT_ID,required"`
	ClientSecret string   `env:"SPOTIFY_CLIENT_SECRET,required"`
	BaseURL      *url.URL `env:"SPOTIFY_PROXY_BASE_URI" envDefault:"http://localhost:27228"`
	APIKey       string   `env:"SPOTIFY_PROXY_API_KEY"`
}

func main() {
	envCfg = new(EnvConfig)
	if err := env.Parse(envCfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	if envCfg.APIKey == "" {
		envCfg.APIKey = randString()
		fmt.Println("APIKey:  ", envCfg.APIKey)
	}

	token = randString()

	authUrl := *envCfg.BaseURL
	authUrl.Path = path.Join(authUrl.Path, "authorize")
	authUrl.RawQuery = url.Values{"token": {token}}.Encode()

	fmt.Println("Auth URL:", authUrl.String())

	redirectUrl := *envCfg.BaseURL
	redirectUrl.Path = path.Join(redirectUrl.Path, "spotify_callback")

	config = &oauth2.Config{
		ClientID:     envCfg.ClientID,
		ClientSecret: envCfg.ClientSecret,
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
	http.HandleFunc("/api/v1/token/terraform", APIToken)
	http.HandleFunc("/spotify_callback", SpotifyCallback)

	log.Fatal(http.ListenAndServe(":27228", nil))
}

// APIToken is the endpoint for refreshing API tokens
func APIToken(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != id || password != envCfg.APIKey {
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

	state := randString()
	mac := hmac.New(sha256.New, []byte(envCfg.APIKey))
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
	if len(state) < 64 {
		fmt.Fprintln(w, "Could not complete authorization: invalid state value")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	state1 := []byte(state[:64])
	stateMac, err := base64.RawURLEncoding.DecodeString(state[64:])
	if err != nil {
		fmt.Fprintln(w, "Could not complete authorization: invalid state value")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mac := hmac.New(sha256.New, []byte(envCfg.APIKey))
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

func randString() string {
	output := make([]byte, 48)
	n, err := rand.Reader.Read(output)
	if err != nil {
		panic(err)
	}
	if n != 48 {
		panic("could not read 48 bytes")
	}
	return base64.RawURLEncoding.EncodeToString(output)
}
