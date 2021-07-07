# syntax=docker/dockerfile:1
# FROM golang:1.16
# WORKDIR /go/src/github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
# RUN go get -u github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy  
# COPY app.go .
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# FROM alpine:latest  
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=0 /go/src/github.com/alexellis/href-counter/app .
# RUN go get -u github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
# # RUN ls
# WORKDIR $GOPATH/src/github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
# RUN go get -u github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
# RUN go install -v ./...
# CMD ["spotify_auth_proxy"]  


FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY . /src/
RUN go get -u github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
RUN CGO_ENABLED=0 go build -o /bin/spotify_auth_proxy

FROM alpine:latest
COPY --from=build /bin/spotify_auth_proxy /bin/spotify_auth_proxy
ENTRYPOINT ["/bin/spotify_auth_proxy"]
