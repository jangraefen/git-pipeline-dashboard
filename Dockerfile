# =============================================================================
# BUILD CONTAINER
# =============================================================================
FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS build

# Create application directory
RUN mkdir /app
ADD . /app/
WORKDIR /app

# Build the application
RUN CGO_ENABLED=0 go mod download
RUN CGO_ENABLED=0 go build -o run .

# =============================================================================
# RUNTIME CONTAINER
# =============================================================================
FROM --platform=$BUILDPLATFORM alpine:latest

# Create application directory
RUN mkdir /app
WORKDIR /app

# Add the execution user
RUN adduser -S -D -H -h /app execuser
USER execuser

# Copy binary from build container
COPY --from=build /app/run /app/run

# Run the application
ENTRYPOINT ["./run"]
