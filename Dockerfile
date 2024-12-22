############################
# STEP 1 build executable binary
############################

FROM golang:1.23-alpine AS builder

# Install git & SSL ca-certificates
# Git is required for fetching the dependencies
# Ca-certificates are required to call HTTPS endpoints
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Copy app for building app binary
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/product-service

############################
# STEP 2 build a small image
############################
FROM alpine:latest
RUN apk update && apk add --no-cache postgresql-client curl

# Create unprivileged appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Make image identifiable
ARG git_commit=unspecified
LABEL git_commit=$git_commit
ARG git_branch=unspecified
LABEL git_branch=$git_branch

RUN mkdir /app
WORKDIR /app

# Copy our static executable
COPY --from=builder /app/bin/product-service .
RUN chown -R appuser:appuser /app

# Use an unprivileged user
USER appuser:appuser

ARG APP_BUILD
ENV APP_BUILD=$APP_BUILD
ARG APP_DEPLOYED_AT
ENV APP_DEPLOYED_AT=$APP_DEPLOYED_AT

# Run the binary
ENTRYPOINT ["/app/product-service"]