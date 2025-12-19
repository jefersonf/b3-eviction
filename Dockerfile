FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /voting-app

COPY . .

RUN go mod init b3e

RUN go mod tidy

RUN go build -ldflags="-s -w" -o /voting-service ./cmd/api/

FROM alpine:3.21 AS final

COPY --from=builder /voting-service /bin/voting-service

EXPOSE 8080

CMD [ "bin/voting-service" ]
