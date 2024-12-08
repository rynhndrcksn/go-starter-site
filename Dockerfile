FROM golang:1.23-bullseye AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64_web ./cmd/web

FROM scratch
WORKDIR /app
COPY --from=build /app/bin/linux_amd64_web ./web
EXPOSE 4000
CMD ["/app/web"]
