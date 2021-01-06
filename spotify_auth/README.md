# spotify_auth

Simple program to return a spotify auth code via [Authorization Code Flow with Proof Key for Code Exchange](https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow-with-proof-key-for-code-exchange-pkce).

## Installation

`go get -u github.com/conradludgate/terraform-provider-spotify/spotify_auth`

## Usage

Run `spotify_auth`.

If the program is able to open a browser tab

1)  find the newly created tab
2)  authorize the application
3)  Close the tab and refer back to the terminal
4)  Take note of the returned `auth code` and `code verifier`

Otherwise

1)  Copy the URL presented
2)  Open the URL in a browser
3)  authorize the application
4)  Copy the URL of the redirect and refer back to the terminal
5)  Paste the URL into the terminal
6)  Take note of the returned `auth code` and `code verifier`
