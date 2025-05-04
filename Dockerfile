FROM golang:1.24.2 AS builder

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
RUN git config --global http.sslVerify false
ENV GOPROXY=direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod download
COPY . .

RUN go build -o main .

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=builder /app/main /app/main

EXPOSE 8080

ENTRYPOINT ["/app/main"]