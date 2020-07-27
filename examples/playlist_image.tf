# provider "local" {
#     version = "1.4.0"
# }

# data "local_file" "tf_logo" {
#     filename = "${path.module}/images/tf.jpg"
# }

# resource "spotify_playlist_image" "playlist_image" {
#     playlist_id = spotify_playlist.playlist.id
#     image_data = data.local_file.tf_logo.content_base64
# }