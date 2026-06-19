# Multi-stage compilation using the latest official Go Alpine image
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

# Inject native platform compilation variables from the build process
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Leverage caching by copying mod definitions first
COPY . .

# Build a hardened, static binary without CGO dependency bindings
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w" \
    -o /app/main main.go

# --- Final Phase Runtime Footprint ---
FROM scratch

# Copy compiled binary out of the builder footprint
COPY --from=builder /app/main /main

EXPOSE 5555

ENTRYPOINT ["/main"]
