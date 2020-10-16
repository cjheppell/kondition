FROM golang:1.15 AS builder

RUN mkdir -p /kondition/

ADD . /kondition/

WORKDIR /kondition/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /kondition/kondition .

FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /kondition/kondition /kondition

RUN chmod 777 /kondition

# Run this container as non-root user
USER 5000

ENTRYPOINT ["/kondition"]
CMD []