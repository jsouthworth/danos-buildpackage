# docker build -t danos-1908-build -f Dockerfile .
# docker run --rm -v $PWD:/mnt/src -v $PWD:/mnt/output  danos-1908-build
FROM debian:buster

RUN mkdir -p '/mnt/src' && \
    mkdir -p '/mnt/output' && \
    mkdir -p '/mnt/pkgs' && \
    mkdir -p /build && \
    groupadd -g 1000 builduser && \
    useradd -r -u 1000 -g builduser -d /home/builduser -m builduser && \
    sed 's/main$/main contrib/' /etc/apt/sources.list && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get -y install build-essential devscripts wget && \
    wget http://ftp.us.debian.org/debian/pool/main/d/devscripts/devscripts_2.20.3~bpo10+1_amd64.deb && \
    dpkg -i devscripts_2.20.3~bpo10+1_amd64.deb && \
    rm devscripts_2.20.3~bpo10+1_amd64.deb

COPY buildpackage /usr/local/bin

WORKDIR /build/src

ENTRYPOINT /usr/local/bin/buildpackage
