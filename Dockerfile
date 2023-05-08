FROM golang:1.18-alpine AS build

WORKDIR /src

COPY main.go .
COPY constants/ constants/
COPY cmd/ cmd/
COPY pkg/ pkg/

COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /prom-logstash-exporter .

FROM alpine:latest AS final

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build /prom-logstash-exporter ./prom-logstash-exporter

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 2112

ENTRYPOINT ["./prom-logstash-exporter"]
