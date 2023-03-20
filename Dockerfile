FROM golang:alpine as builder
WORKDIR /goscript
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goscript -ldflags="-w -s" .

FROM golang:alpine

WORKDIR /app
COPY --from=builder /goscript/goscript /app/goscript
ENTRYPOINT ["/app/goscript"]