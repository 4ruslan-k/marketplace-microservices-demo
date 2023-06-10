FROM golang:1.20.5-alpine

RUN apk add --no-cache git

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /go/bin/app cmd/app/main.go

EXPOSE 4001

CMD ["/go/bin/app"]