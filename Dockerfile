FROM golang:alpine AS golang-builder
WORKDIR /build
COPY . .
RUN apk add build-base
RUN go mod download
RUN go build ./cmd/microurl
WORKDIR /dist
RUN cp /build/microurl .

FROM alpine
EXPOSE 3000
COPY --from=golang-builder /build/microurl /
COPY ./config.json /config.json
COPY ./banner /banner
RUN mkdir /data
ENV CONFIG_TLSCRT=/data/server.crt
ENV CONFIG_TLSKEY=/data/server.key
ENV CONFIG_TLSENABLED=true

ENTRYPOINT ["/microurl"]

