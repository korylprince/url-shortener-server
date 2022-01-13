FROM golang:1-alpine as builder

ARG VERSION

RUN go install "github.com/korylprince/url-shortener-server/v2@$VERSION"

FROM alpine:3.15

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/url-shortener-server /shortener

CMD ["/shortener"]
