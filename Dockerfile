FROM golang:1.20.2

WORKDIR /go/src/GeekHome

COPY . .

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=$GOPROXY

ENV ENVIRONMENT=Prod

RUN go build -o main ./src/main

EXPOSE 5466

CMD ["./main"]