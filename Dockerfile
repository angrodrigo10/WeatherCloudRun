FROM golang:1.23.2 as builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates
RUN GOOS=linux CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server

CMD ["/server"]
