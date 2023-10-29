FROM golang:1.21-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod  \
    --mount=type=cache,target=/root/.cache/go-build go mod download  \
    && go build -o main ./internal/simple-app/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
