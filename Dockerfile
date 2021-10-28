FROM golang:1.17.1-alpine3.14 AS phprom-builder

COPY ./ /go/src/github.com/chaseisabelle/phprom

WORKDIR /go/src/github.com/chaseisabelle/phprom/cmd/v1

RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /phprom


FROM alpine:3.14 AS phprom

COPY --from=phprom-builder /phprom /phprom

CMD ["/phprom"]
