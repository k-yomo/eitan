FROM golang:1.16
ENV GO111MODULE=on

WORKDIR /go/src/github.com/k-yomo/eitan

COPY go.mod go.sum ./
RUN go mod download

RUN go get -u github.com/cosmtrek/air

COPY src/eitan_service ./src/notification_service
COPY src/internal ./src/internal
COPY src/pkg ./src/pkg

CMD cd src/notification_service && air

