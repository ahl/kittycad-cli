FROM golang:alpine as builder
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/kittycad/cli

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/kittycad/cli \
	&& make static \
	&& mv cli /usr/bin/cli \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/cli /usr/bin/cli
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "cli" ]
CMD [ "--help" ]
