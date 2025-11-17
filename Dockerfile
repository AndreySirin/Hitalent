FROM golang:1.24.7 AS builder

WORKDIR /app

COPY . ./

RUN go mod tidy

WORKDIR /app/cmd

RUN go build -o /app/hitalent


FROM ubuntu:latest

WORKDIR /root

COPY --from=builder /app/hitalent .

EXPOSE 8080

CMD ["/root/hitalent"]