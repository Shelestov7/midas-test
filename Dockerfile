FROM golang:1.18 as builder

WORKDIR /midas

COPY . ./

RUN go build -o /app

FROM ubuntu:latest

WORKDIR /

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

COPY --from=builder /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
