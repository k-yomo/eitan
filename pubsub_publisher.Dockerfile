# build
FROM golang:1.16-alpine as build
ENV GO111MODULE=on

WORKDIR /go/src/app

RUN apk --no-cache add make ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY src/pubsub_publisher src/pubsub_publisher
COPY src/internal src/internal
COPY src/pkg src/pkg

RUN CGO_ENABLED=0 go build -o bin/server -ldflags "-w -s" ./src/pubsub_publisher

# exec
FROM scratch
COPY --from=build /go/src/app/bin/server ./server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
ENTRYPOINT ["./server"]
