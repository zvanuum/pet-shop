FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go get github.com/dghubble/oauth1
RUN go get github.com/sajari/fuzzy
RUN go get github.com/spf13/viper
RUN go build -o pet-shop .

CMD ["/app/pet-shop"]