# terraform-provider-spotify

This is a terraform provider for managing your spotify playlists.

## Installation

To install the provider, run
```sh
cd tf
make install
```
This will build the binary and install it into the terraform plugins dir.

## How to use

Currently this provider requires 2 environment variables:
`SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET`.
You can get these by going to https://developer.spotify.com/dashboard and registering an application.
On the application, you must have `http://localhost:27228/spotify_callback` as a registered redirect URL.

When you run `terraform plan` or `terraform apply`, it will open the spotify login page in your browser.
When you approve, you will be redirected to `localhost:27228` which the provider will be listening on.
(If JS is enabled, the tab should auto close. This is convenient when you've already approved since spotify will auto redirect.)

I may instead make it use [implicit grant authentication](https://developer.spotify.com/documentation/general/guides/authorization-guide/#implicit-grant-flow)
since it doesn't require a secret key, and then I could release this with my client-id.

## Example

```tf
# Creates a private playlist named "My playlist"
# and adds 2 tracks to it
resource "spotify_playlist" "playlist" {
    name = "My playlist"
    description = "My playlist is so awesome"
    public = false

    tracks = [
        data.spotify_search_track.overkill.id,
        data.spotify_search_track.blackwater.id,
    ]
}

# Searches spotify for "Overkill RIOT", returns the first track ID
data "spotify_search_track" "overkill" {
    name = "Overkill"

    artists = [
        "RIOT",
    ]
}

# Searches spotify for "Blackwater RIOT", returns the first track ID
data "spotify_search_track" "blackwater" {
    name = "Blackwater"

    artists = [
        "RIOT",
    ]
}
```