FROM golang:1.14

ENV CGO_ENABLED=0\
    GOOS=linux

WORKDIR /go/src/github.com/wincus/k8s-event-logger

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netcgo -ldflags '-w' -o k8s-event-logger main.go

FROM alpine:latest

WORKDIR /

COPY --from=0 /go/src/github.com/wincus/k8s-event-logger/k8s-event-logger .

ENTRYPOINT ["/k8s-event-logger"]