FROM golang:1.7

RUN mkdir -p /app

WORKDIR /app

COPY . /app

RUN go build

CMD ["./app"]
