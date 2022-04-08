FROM golang:1.12.0-alpine3.9

RUN apk update

RUN apk add git

RUN mkdir /app

ADD ./app /app

WORKDIR /app

RUN go mod download

RUN go build -o main .

WORKDIR /src

CMD ["/app/main"]