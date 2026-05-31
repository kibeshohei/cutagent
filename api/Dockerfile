# Build context: repository root (docker build -f api/Dockerfile .)
FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY api/go.mod api/go.sum ./api/
RUN cd api && go mod download
COPY api ./api
RUN cd api && CGO_ENABLED=0 go build -o /server ./cmd/server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /server /server
COPY prompts /app/prompts
ENV PROMPTS_DIR=/app/prompts
EXPOSE 8080
ENTRYPOINT ["/server"]
