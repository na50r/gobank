
FROM golang:1.23.10-alpine AS builder
RUN apk add --no-cache gcc g++ git openssh-client
RUN mkdir /build
RUN mkdir /build/sse
COPY go.mod go.sum storage.go types.go api.go main.go /build/
COPY sse/sse.go /build/sse/
WORKDIR /build
# Workaround to not having to change the imports in the monorepo
RUN sed -i 's|github.com/na50r/gobank/backend/sse|github.com/na50r/build/sse|' api.go
RUN sed -i 's|github.com/na50r/gobank/backend/sse|github.com/na50r/build/sse|' main.go
RUN sed -i 's|module github.com/na50r/gobank/backend|module github.com/na50r/build|' go.mod
RUN go mod tidy

# Required to get sqlite3 to work
RUN CGO_ENABLED=1 go build -o bin/gobank

FROM alpine
# Reference1: https://dev.to/heroku/deploying-your-first-golang-webapp-11b3
# Reference2: https://stackoverflow.com/questions/45972608/how-to-give-folder-permissions-inside-a-docker-container-folder
RUN adduser -S -D -H -h /app appuser
COPY --from=builder /build/bin/gobank /app/
WORKDIR /app
RUN chown -R appuser /app && chmod 755 /app
USER appuser
CMD ["./gobank"]

