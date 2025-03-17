FROM alpine:3.19

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file iproute2 iputils-ping

VOLUME /cometbft
WORKDIR /cometbft
EXPOSE 26656 26657 26660
ENTRYPOINT ["/usr/bin/wrapper.sh"]

STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh
COPY config-template.toml /etc/cometbft/config-template.toml