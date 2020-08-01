# spotify_auth

Simple program to authorize and return a spotify access token via [implicit grant](https://developer.spotify.com/documentation/general/guides/authorization-guide/#implicit-grant-flow).

## Installation

`go install github.com/conradludgate/terraform-provider-spotify/spotify_auth`

## Usage

Run
```sh
spotify_auth | source /dev/stdin
```
And click the link provided.
Login to spotify and authorize the application.
Once spotify redirects, copy the URL and paste it into the terminal.

If everything succeeds, the variable `SPOTIFY_ACCESS_TOKEN` will be set with the access token granted by spotify.

## Configuration

The default values are for [terraform-provider-spotify](https://github.com/conradludgate/terraform-provider-spotify) but all the values are configurable.

*   **SPOTIFY_CLIENT_ID** - Client ID of the spotify application
*   **SPOTIFY_SCOPES** - Scopes required, comma separated (see: https://developer.spotify.com/documentation/general/guides/scopes)
*   **SPOTIFY_REDIRECT_URI** - Redirect URI for the spotify application
