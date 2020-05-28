# docker build -t danos-1908-build -f Dockerfile .
# docker run --rm -v $PWD:/mnt/src -v $PWD:/mnt/output  danos-1908-build
FROM debian:buster

RUN mkdir -p '/mnt/src' && \
    mkdir -p '/mnt/output' && \
    mkdir -p '/mnt/pkgs' && \
    mkdir -p /build/src && \
    groupadd -g 1000 builduser && \
    useradd -r -u 1000 -g builduser -d /home/builduser builduser && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get -y install build-essential devscripts wget

COPY buildpackage /usr/local/bin

WORKDIR /build/src

ENTRYPOINT /usr/local/bin/buildpackage
