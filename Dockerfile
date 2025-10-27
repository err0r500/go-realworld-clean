FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum* ./
# Use BuildKit secrets to configure git credentials without exposing them
# Use cache mount to persist Go modules between builds
RUN go mod download

COPY . .

# Use cache mounts for both Go modules and build cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o app .

FROM scratch
COPY --from=builder /app/app .
EXPOSE 8080
ENTRYPOINT ["/app"]
