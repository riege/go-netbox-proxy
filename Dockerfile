# Start by building the application.
FROM golang:1.18.3-buster as build

WORKDIR /go/src/go-netbox-proxy
COPY . /go/src/go-netbox-proxy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -o /go/bin/go-netbox-proxy

# Now copy it into our base image.
FROM gcr.io/distroless/static:latest-amd64
COPY --from=build /go/bin/go-netbox-proxy /

LABEL org.opencontainers.image.source https://github.com/riege/go-netbox-proxy

CMD ["/go-netbox-proxy"]
