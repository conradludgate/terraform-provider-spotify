# spotify_auth_proxy

This is an instance of a 'spotify auth server' which acts as an interface between a client and the spotify oauth API.

## Installation

With `go` installed, run

```sh
go get -u github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
```

## Usage

First, you need a Spotify client id and secret. Visit https://developer.spotify.com/dashboard/ to create an application.

If you plan to run this proxy locally, set the redirect URI of the application to `http://localhost:27228/spotify_callback`.

```sh
export SPOTIFY_CLIENT_REDIRECT_URI=http://localhost:27228/spotify_callback
```

You will also need to register the callback URI with Spotify for your application. Visit https://developer.spotify.com/dashboard/, click on your application, find and click the "Edit Settings" button, and paste the `spotify_callback` URI above into "Redirect URIs". Scroll down and click "Save".

If you plan to host it on an external server, the redirect URI should be the equivalent URL path on your host. I would recommend putting it behind Nginx or some other reverse proxy with SSL enabled.

To start the server, make sure the environment variables are set `SPOTIFY_CLIENT_ID`, `SPOTIFY_CLIENT_SECRET` and `SPOTIFY_CLIENT_REDIRECT_URI`, then run

```sh
spotify_auth_proxy
```

It should output the following:

```
APIKey: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
Token:  YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY
```

Take note of both of these values.

Now, open a browser and navigate to `http://localhost:27228/authorize?token={TOKEN FROM ABOVE}` or approriate for your server location. It should redirect you to spotify to login, and then you will be redirected back to the page where it should confirm that you've authorized correctly.

The API Key is how you will retrieve the access token. The server will handle the token expiration and refreshes for you.
