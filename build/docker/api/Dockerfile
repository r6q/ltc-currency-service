FROM golang:1.15-alpine

COPY . /tmp/api_app
WORKDIR /tmp/api_app

RUN go build -o app ./cmd/api

RUN mv app config /go
WORKDIR /go
RUN rm -rf /tmp/api_app

RUN apk add --no-cache openssl
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENTRYPOINT ["dockerize", "-wait", "tcp://database:3306", "-timeout", "10s", "./app", "-config=config/application-docker.yaml"]
