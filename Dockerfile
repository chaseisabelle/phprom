FROM golang:latest AS phprom

RUN go get github.com/chaseisabelle/phprom/cmd/v1

WORKDIR /go/src/github.com/chaseisabelle/phprom/cmd/v1

RUN go get -v && go build -o phprom

CMD ["./phprom"]
