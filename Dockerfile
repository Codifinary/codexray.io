FROM golang:1.23-bullseye AS backend-builder
RUN apt update && apt install -y liblz4-dev
# Set Go Environment
ENV GOARCH=amd64
ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /tmp/src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG VERSION=unknown
RUN go build -mod=readonly -ldflags "-X main.version=$VERSION" -o codexray .


FROM debian:bullseye
RUN apt update && apt install -y ca-certificates && apt clean

WORKDIR /opt/codexray
COPY --from=backend-builder /tmp/src/codexray /opt/codexray/codexray

VOLUME /data
EXPOSE 8080

ENTRYPOINT ["/opt/codexray/codexray"]
