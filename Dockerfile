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
    echo "deb http://deb.debian.org/debian buster-backports main" > \
	/etc/apt/sources.list.d/backports.list && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get -y install devscripts git-buildpackage \
        pristine-tar bzip2 xz-utils wget && \
    apt-get -y -t buster-backports install devscripts

COPY buildpackage /usr/local/bin

WORKDIR /build

ENTRYPOINT /usr/local/bin/buildpackage
