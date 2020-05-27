FROM golang
ADD ./web_server /web_server
ADD ./secret.go /web_server
ADD ./secureConfig.json /web_server

WORKDIR /web_server

RUN go build

ENTRYPOINT ./web_server

EXPOSE 8000