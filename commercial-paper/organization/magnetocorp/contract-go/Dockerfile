FROM golang:1.15-buster as setup

ENV GO111MODULE=on
WORKDIR /go/src/contract-go/
ADD ./go.* ./

RUN go mod download

COPY ./commercial-paper/ ./commercial-paper/
COPY ./ledger-api/ ./ledger-api/
COPY ./main.go .

FROM setup as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /go/bin/contract-cc /go/src/contract-go/.


# FROM gcr.io/distroless/base-debian10
FROM alpine:3.12

COPY --from=builder /go/bin/contract-cc /
ENTRYPOINT ["/contract-cc"] 