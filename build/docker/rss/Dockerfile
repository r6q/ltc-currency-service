FROM golang:1.15-alpine

COPY . /tmp/rss_app
RUN pwd
WORKDIR /tmp/rss_app

RUN go build -o app ./cmd/rss

RUN mv app config /go
WORKDIR /go
RUN rm -rf /tmp/rss_app

ENTRYPOINT ["./app", "-config=config/application-docker.yaml"]
