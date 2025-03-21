FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build_linux

FROM scratch

COPY --from=builder /app/build/proxier /proxier

EXPOSE 8080

ENTRYPOINT ["/proxier", "-config", "/config.yaml"]
