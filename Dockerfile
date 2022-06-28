FROM golang:1.18-alpine AS builder
RUN apk add --update ca-certificates

WORKDIR /go/src/github.com/broothie/dillbook
COPY . .
RUN go build cmd/server/main.go

FROM alpine:3.16
COPY --from=builder /go/src/github.com/broothie/dillbook/main main

CMD ["./main"]
