# docker build -t danos-2015-build -f Dockerfile .
# docker run --rm -v $PWD:/mnt/src -v $PWD:/mnt/output  danos-2015-build
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
    apt-get -y install build-essential devscripts wget \
        git-buildpackage pristine-tar bzip2 xz-utils && \
    apt-get -y -t buster-backports install devscripts && \
    echo "deb http://s3-us-west-1.amazonaws.com/2105.repos.danosproject.org/repo/ 2105 main" > /etc/apt/sources.list.d/danos.list && \
    wget -q -O- https://s3-us-west-1.amazonaws.com/repos.danosproject.org/Release.key | apt-key add - && \
    apt-get update

COPY buildpackage /usr/local/bin

WORKDIR /build

ENTRYPOINT /usr/local/bin/buildpackage
