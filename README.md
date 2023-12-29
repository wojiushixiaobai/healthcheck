# Docker healthchecks in distroless

This is a simple example of how to use healthchecks in a Dockerfile using distroless images.

cat Dockerfile
```yaml
# Start by building the application.
FROM golang:1.18 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /

# download the healthcheck binary
ADD check /usr/bin/check
CMD ["/app"]
```

cat docker-compose.yml
```yaml
version: "3.9"
services:
  app:
    build: .
    image: examples/app
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "check", "http://localhost:8080/health"]
      interval: 3s
      timeout: 3s
      retries: 10
```