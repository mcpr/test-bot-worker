FROM golang:1.8


WORKDIR /go/src/test-bot-worker

COPY . .

RUN go get github.com/sparrc/gdm && \
    gdm restore && \
    go build -v && \
    mv test-bot-worker /usr/local/bin