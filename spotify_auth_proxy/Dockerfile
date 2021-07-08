FROM golang:1.16-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /bin/spotify_auth_proxy

FROM alpine:latest
COPY --from=build /bin/spotify_auth_proxy /bin/spotify_auth_proxy
ENTRYPOINT ["/bin/spotify_auth_proxy"]
