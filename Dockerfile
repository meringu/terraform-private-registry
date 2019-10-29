FROM golang:1.13-buster AS golang

ADD . /src
WORKDIR /src
RUN go build \
    -a \
    -o /terraform-private-registry \
    -ldflags '-linkmode external -extldflags -s -w' \
    main.go

FROM alpine

COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=golang /terraform-private-registry /

ENTRYPOINT ["/terraform-private-registry"]
