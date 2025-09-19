FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /narde-server ./cmd/service/

FROM alpine:latest

WORKDIR /

COPY --from=builder /narde-server /narde-server

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/narde-server"]