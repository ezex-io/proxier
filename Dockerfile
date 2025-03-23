# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download -x
RUN make build_linux

FROM alpine:latest

RUN mkdir /etc/proxier
COPY --from=builder /app/build/proxier /usr/bin/proxier
COPY --from=builder /app/docs/config.example.yml /etc/proxier/config.yml

EXPOSE 8080

ENTRYPOINT ["/usr/bin/proxier", "-config", "/etc/proxier/config.yml"]
