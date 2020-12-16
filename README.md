# danos-buildpackage

This tool aids in building DANOS packages using containers. It uses the docker API and should work with any compatibile container engine.

## Installation

Binaries are included in the Releases section of this repo, the appropriate binary for your operating system can be placed in your PATH.

From source:
```
$ git clone https://github.com/jsouthworth/danos-buildpackage
$ cd danos-buildpackage
$ go install jsouthworth.net/go/danos-buildpackage/cmd/danos-buildpackage
```

## Usage

 - Clone the desired repository.
 - "cd <the-cloned-source-tree>"
 - Checkout the correct tag for the desired release
   - "git checkout danos/<version>"
 - Run:
   - "danos-buildpackage -version <version>"

The debian files from completed build will be in the parent directory.

Note: Some packages require ipv6 support to build properly. To enable ipv6 support in docker add the following to "/etc/docker/daemon.json" or the equivalent on your platform:
```
{
	"ipv6": true,
	"fixed-cidr-v6": "2001:db8:1::/64"
}
```
The "2001:db8:1::/64" address is for documenation only, one should replace the "fixed-cidr-v6" address with one's own subnet or a [ULA](https://tools.ietf.org/html/rfc4193).

The packages that currently require ipv6 to build are:
 - vyatta-vrrp
 - strongswan

### Advanced options
```
$ danos-buildpackage -h
Usage of danos-buildpackage:
  -dest string
    	destination directory (default "..")
  -image-name string
    	name of docker image (default "jsouthworth/danos-buildpackage")
  -local
    	is the image only on the local system
  -no-clean
    	don't delete the container when done
  -pkg string
    	preferred package directory
  -src string
    	source directory (default ".")
  -version string
    	version of danos to build for (default "latest")
```
