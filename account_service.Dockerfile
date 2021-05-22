# build
FROM golang:1.16-alpine as build
ENV GO111MODULE=on

WORKDIR /go/src/app

RUN apk --no-cache add make ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.2 && \
    wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY src/account_service src/account_service
COPY src/internal src/internal
COPY src/pkg src/pkg

RUN CGO_ENABLED=0 go build -o bin/server -ldflags "-w -s" ./src/account_service

# exec
FROM scratch
COPY --from=build /go/src/app/bin/server ./server
COPY --from=build /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT ["./server"]