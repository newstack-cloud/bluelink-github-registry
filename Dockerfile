###################
# BUILD APPLICATION BINARY
###################

# Please keep up to date with the new-version of Golang docker for builder
FROM golang:1.24.2-bookworm AS build

# Add non-root user
RUN useradd -u 1001 bluelink

RUN apt update && \
  apt upgrade -y && \
  apt install -y git build-essential curl && \
  apt install openssh-server -y

WORKDIR /app/bluelink_github_registry

# Copy go module files to load dependencies.
COPY go.mod ./go.mod
COPY go.sum ./go.sum

COPY cmd ./cmd
COPY internal ./internal

RUN go mod download

RUN go build \
  # Build with static linking to ensure everything is included in the binary
  # which allows us to run the binary with scratch.
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o bluelink_github_registry \
  cmd/main.go

###################
# PRODUCTION IMAGE (lean image to run pre-built binary)
###################

FROM scratch

WORKDIR /

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/bluelink_github_registry/bluelink_github_registry /bluelink_github_registry

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Use non-root user
USER bluelink

EXPOSE 8085

CMD ["/bluelink_github_registry"]