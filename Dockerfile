FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression. Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

FROM golang:1.14-stretch AS builder

ARG APP_VERSION=Unknown
ARG APP_COMMIT=Unknown
ARG APP_CREATED=Unknown

# Services:
RUN apt-get update && apt-get install -y openssh-client
RUN GO111MODULE=off go get -u github.com/gobuffalo/packr/v2/packr2

# Preparing project env:
RUN mkdir -m 0777 -p /go/src/github.com/spyzhov/healthy

WORKDIR /go/src/github.com/spyzhov/healthy

ENV GOOS=linux
ENV GOARCH=amd64

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go mod vendor && \
    packr2 && \
	go build \
	    -mod=vendor \
	    -ldflags "\
	        -w \
	        -X github.com/spyzhov/healthy/app.Version=${APP_VERSION} \
	        -X github.com/spyzhov/healthy/app.Commit=${APP_COMMIT} \
	        -X github.com/spyzhov/healthy/app.Created=${APP_CREATED} \
	    " \
	    -o /go/bin/healthy .

FROM debian:9-slim
# environment
ENV LOG_LEVEL=info
ENV PORT=80
ENV MANAGEMENT_PORT=3280
# configurations
EXPOSE 80
EXPOSE 3280
WORKDIR /root
# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# the main program:
COPY --from=builder /go/bin/healthy ./healthy
COPY --from=builder /go/src/github.com/spyzhov/healthy/example.yaml ./example.yaml
CMD ["./healthy"]
