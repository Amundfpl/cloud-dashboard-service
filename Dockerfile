# ---------- STAGE 1: Builder ----------
FROM golang:1.23 AS builder

LABEL maintainer="mail@host.tld"
LABEL stage=builder

# Set working directory
WORKDIR /go/src/app

# Copy go mod files first (caching)
COPY go.mod go.sum ./
RUN go mod download


COPY cmd ./cmd
COPY cache ./cache
COPY db ./db
COPY handlers ./handlers
COPY httpclient ./httpclient
COPY server ./server
COPY services ./services
COPY utils ./utils
COPY static ./static

# Set working directory to where main.go is
WORKDIR /go/src/app/cmd

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /server

# ---------- STAGE 2: Runtime ----------
FROM alpine:latest

# Install certificates (for HTTPS support)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder
COPY --from=builder /server /server

COPY --from=builder /go/src/app/static /static

# Expose the port your app listens on (change if not 8080)
EXPOSE 8080

# Start the binary
ENTRYPOINT ["/server"]

