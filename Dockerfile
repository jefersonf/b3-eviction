FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /voting-app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /voting-service ./cmd/api/

FROM alpine:3.21 AS final

COPY --from=builder /voting-service /bin/voting-service

EXPOSE 8080

CMD [ "bin/voting-service" ]