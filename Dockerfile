FROM --platform=$BUILDPLATFORM golang:latest AS build
ARG TARGETOS
ARG TARGETARCH
COPY . /otpauth
WORKDIR /otpauth
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN go build

FROM scratch
COPY --from=build /otpauth/otpauth /app/otpauth
ENTRYPOINT ["/app/otpauth"]
