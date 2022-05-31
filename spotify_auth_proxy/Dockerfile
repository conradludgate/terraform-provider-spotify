FROM golang:1.16-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o /bin/spotify_auth_proxy

FROM alpine:latest

WORKDIR /home/spotify_auth_proxy
RUN addgroup -S spotify && \
    adduser -S spotify_auth_proxy -G spotify

USER spotify_auth_proxy
COPY --from=build /bin/spotify_auth_proxy .local/bin/spotify_auth_proxy
ENTRYPOINT [".local/bin/spotify_auth_proxy"]
