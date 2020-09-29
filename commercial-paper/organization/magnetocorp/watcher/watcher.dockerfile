FROM golang:1.15-buster as setup

ENV GO111MODULE=on

WORKDIR /go/src/watcher/
ADD ./go.* ./

RUN go mod download

COPY ./main.go .

FROM setup as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /go/bin/watcher /go/src/watcher/.

FROM alpine:3.12

COPY --from=builder /go/bin/watcher /

# CMD [ "/bin/sh", "-c", "sleep 3000" ]
ENTRYPOINT ["/watcher"]
