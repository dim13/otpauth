FROM golang:latest AS build
COPY . /otpauth
WORKDIR /otpauth
ENV CGO_ENABLED=0
RUN go build

FROM scratch
COPY --from=build /otpauth/otpauth /app/otpauth
ENTRYPOINT ["/app/otpauth"]
