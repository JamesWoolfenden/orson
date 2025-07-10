FROM alpine:3.18.2

RUN apk --no-cache add build-base git curl jq bash
RUN curl -s -k https://api.github.com/repos/JamesWoolfenden/orson/releases/latest | jq '.assets[] | select(.name | contains("linux_386")) | select(.content_type | contains("gzip")) | .browser_download_url' -r | awk '{print "curl -L -k " $0 " -o ./orson.tar.gz"}' | sh
RUN tar -xf ./orson.tar.gz -C /usr/bin/ && rm ./orson.tar.gz && chmod +x /usr/bin/orson && echo 'alias orson="/usr/bin/orson"' >> ~/.bashrc
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]

LABEL layer.0.author="JamesWoolfenden"
