FROM golang:1.9.2 as go-builder
WORKDIR /go
RUN go get github.com/kelseyhightower/envconfig
RUN go get github.com/gtaylor/factorio-rcon
RUN go get golang.org/x/crypto/acme/autocert
COPY *.go /go/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
RUN mkdir /app && cp main /app/main

FROM scratch
COPY --from=go-builder /app /app
COPY public /app/public
WORKDIR /app
CMD ["/app/main"]
