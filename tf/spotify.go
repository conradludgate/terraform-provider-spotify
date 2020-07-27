package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/skratchdot/open-golang/open"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var spotifyClient *spotify.Client

func ClientConfigurer(_ *schema.ResourceData) (interface{}, error) {
	if spotifyClient != nil {
		return spotifyClient, nil
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: "http://localhost:27228/spotify_callback",
		Scopes: []string{
			"user-read-email",
			"user-read-private",
			"playlist-read-collaborative",
			"playlist-read-private",
			"playlist-modify-private",
			"playlist-modify-public",
			"user-library-read",
			"user-library-modify",
			"ugc-image-upload",
		},
	}

	httpClient, err := auth(config)
	if err != nil {
		return nil, err
	}

	client := spotify.NewClient(httpClient)
	spotifyClient = &client
	return spotifyClient, nil
}

func auth(config *oauth2.Config) (*http.Client, error) {
	ctx := context.Background()

	state := uuid.New().String()

	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	if err := open.Run(url); err != nil {
		return nil, err
	}

	var code string
	authServer := &http.Server{Addr: ":27228"}
	http.HandleFunc("/spotify_callback", func(w http.ResponseWriter, r *http.Request) {
		if state != r.FormValue("state") {
			w.Write([]byte("not accepted"))
			return
		}

		code = r.FormValue("code")

		w.Write([]byte("<body>accepted. you can now close this tab<script>window.close()</script></body>"))

		go authServer.Shutdown(ctx)
	})

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := authServer.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	wg.Wait()

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return config.Client(ctx, token), nil
}
