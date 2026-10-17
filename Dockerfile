FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1

RUN go build -ldflags "-linkmode external -extldflags -static" \
    -o /app/server ./cmd/server

RUN mkdir -p /app/plugins && \
    for d in builtins/*/; do \
        go build -buildmode=plugin \
            -ldflags "-linkmode external -extldflags -static" \
            -o "/app/plugins/$(basename "$d").so" "./$d"; \
    done

FROM alpine:latest AS certs

FROM scratch

COPY --from=builder /app/ /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/server"]
