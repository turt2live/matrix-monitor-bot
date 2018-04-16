FROM docker.io/alpine
COPY . /tmp/src
RUN apk add --no-cache \
      su-exec \
 && apk add --no-cache \
      -t build-deps \
      go \
      git \
      musl-dev \
      dos2unix \
 && apk add --no-cache ca-certificates \
 && cd /tmp/src \
 && GOPATH=`pwd` go get github.com/constabulary/gb/... \
 && PATH=$PATH:`pwd`/bin gb vendor restore \
 && GOPATH=`pwd`:`pwd`/vendor go build -o bin/monitor_bot ./src/github.com/turt2live/matrix-monitor-bot/cmd/ \
 && cp bin/monitor_bot .docker/run.sh /usr/local/bin \
 && cp config.sample.yaml /etc/monitor-bot.yaml.sample \
 && dos2unix /etc/monitor-bot.yaml.sample \
 && dos2unix /usr/local/bin/run.sh \
 && cd / \
 && rm -rf /tmp/* \
 && apk del build-deps

CMD /usr/local/bin/run.sh
VOLUME ["/data"]
EXPOSE 8080
