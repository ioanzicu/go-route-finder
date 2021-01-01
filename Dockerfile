# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:latest

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

LABEL maintainer="Ioan Zîcu <ioan.zicu94@gmail.com>"

WORKDIR /route-finder/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . . 

RUN go build

EXPOSE 8000

CMD ["go", "run", "main.go"]
