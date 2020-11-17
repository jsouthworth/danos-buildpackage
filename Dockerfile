# docker build -t danos-2009-build -f Dockerfile .
# docker run --rm -v $PWD:/mnt/src -v $PWD:/mnt/output  danos-2009-build
FROM debian:buster-slim

RUN mkdir -p '/mnt/src' && \
    mkdir -p '/mnt/output' && \
    mkdir -p '/mnt/pkgs' && \
    mkdir -p /build && \
    groupadd -g 1000 builduser && \
    useradd -r -u 1000 -g builduser -d /home/builduser -m builduser && \
    sed 's/main$/main contrib/' /etc/apt/sources.list && \
    echo "deb http://deb.debian.org/debian buster-backports main" > \
	/etc/apt/sources.list.d/backports.list && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get -y install devscripts wget && \
    apt-get -y -t buster-backports install devscripts && \
    echo "deb http://s3-us-west-1.amazonaws.com/2009.repos.danosproject.org/repo/ 2009 main" > /etc/apt/sources.list.d/danos.list && \
    wget -q -O- https://s3-us-west-1.amazonaws.com/repos.danosproject.org/Release.key | apt-key add - && \
    apt-get update

COPY buildpackage /usr/local/bin

WORKDIR /build/src

ENTRYPOINT /usr/local/bin/buildpackage
