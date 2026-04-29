FROM golang:1.22-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /ingestion .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /ingestion /ingestion
ENTRYPOINT ["/ingestion"]
