FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/sneaky ./...

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/sneaky /bin/sneaky
ENTRYPOINT ["/bin/sneaky"]
