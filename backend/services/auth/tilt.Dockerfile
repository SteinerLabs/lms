FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY ./shared ./shared
COPY ./services/auth ./services/auth

WORKDIR /build/services/auth
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/main.go

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/services/auth/service /service
ENTRYPOINT ["/service"]