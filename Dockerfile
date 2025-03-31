# Build stage
FROM golang:1.24-alpine AS builder

RUN apk --no-cache add make git

WORKDIR /app
COPY . .

RUN go mod download -x
RUN make build_linux

FROM alpine:latest

RUN mkdir /etc/proxier
COPY --from=builder /app/build/proxier /usr/bin/proxier

EXPOSE 8080

ENTRYPOINT ["/usr/bin/proxier", "-config", "/etc/proxier/config.yml"]
