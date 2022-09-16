# Build
FROM golang:1.18 AS builder

COPY . /otpauth
WORKDIR /otpauth

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/otpauth

# Run
FROM scratch

COPY --from=builder /otpauth/build/otpauth /app/otpauth

ENTRYPOINT ["/app/otpauth"]
