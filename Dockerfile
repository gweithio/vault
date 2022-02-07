FROM golang:latest

RUN mkdir -p /go/src/vaultApi

WORKDIR /go/src/vaultApi

COPY . /go/src/vaultApi

RUN go install vaultApi

CMD /go/bin/vaultApi

EXPOSE 8080