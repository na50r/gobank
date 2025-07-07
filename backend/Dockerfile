
FROM golang:1.23.10-alpine AS builder
RUN apk add --no-cache gcc g++ git openssh-client
RUN mkdir /build
RUN mkdir /build/sse
COPY go.mod go.sum storage.go types.go api.go main.go /build/
COPY sse/sse.go /build/sse/
WORKDIR /build
RUN sed -i 's|github.com/na50r/gobank/backend/sse|github.com/na50r/build/sse|' api.go
RUN sed -i 's|github.com/na50r/gobank/backend/sse|github.com/na50r/build/sse|' main.go
RUN sed -i 's|module github.com/na50r/gobank/backend|module github.com/na50r/build|' go.mod
RUN go mod tidy
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/gobank


FROM alpine
# Reference1: https://dev.to/heroku/deploying-your-first-golang-webapp-11b3
# Reference2: https://stackoverflow.com/questions/45972608/how-to-give-folder-permissions-inside-a-docker-container-folder
RUN adduser -S -D -H -h /app appuser
COPY --from=builder /build/bin/gobank /app/
WORKDIR /app
RUN chown -R appuser /app && chmod 755 /app
USER appuser
ENV CLIENT="http://localhost:5500"
CMD ["./gobank"]

