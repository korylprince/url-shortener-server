FROM golang:1.14-alpine as builder

ARG VERSION

RUN apk add --no-cache git nodejs npm ca-certificates

RUN git clone --branch "$VERSION" --single-branch --depth 1 --recurse-submodules \
    https://github.com/korylprince/url-shortener-server.git /go/src/github.com/korylprince/url-shortener-server

RUN cd /go/src/github.com/korylprince/url-shortener-server/client && \
    npm install && \
    npm run build-prod

RUN go get -u github.com/gobuffalo/packr/...

RUN /go/bin/packr install github.com/korylprince/url-shortener-server

FROM alpine:3.11

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/url-shortener-server /shortener

CMD ["/shortener"]
