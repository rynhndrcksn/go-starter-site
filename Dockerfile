# Stage 1: Build the application.
FROM golang:1.24-bullseye AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o=./bin/linux_amd64_web ./cmd/web

# Stage 2: Create the final image.
FROM gcr.io/distroless/static-debian12 AS final
WORKDIR /app

# Copy the binary from the build stage.
COPY --from=build --chown=nonroot:nonroot /app/bin/linux_amd64_web ./web

# Switch to the non-root user.
USER nonroot

EXPOSE 4000

# Healthcheck to verify the binary can run and output its version
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/web", "-version"]

CMD ["/app/web"]
