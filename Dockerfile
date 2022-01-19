FROM golang:1.17-alpine

WORKDIR $GOPATH/github.com/nuzurie/shopify

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["shopify"]