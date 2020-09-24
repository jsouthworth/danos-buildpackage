# danos-buildpackage

This tool aids in building DANOS packages using containers. It uses the docker API and should work with any compatibile container engine.

## Installation

Binaries are included for the Release, the appropriate binary for your operating system can be placed in your PATH.

From source:

```
$ git clone https://github.com/jsouthworth/danos-buildpackage
$ cd danos-buildpackage
$ go install jsouthworth.net/go/danos-buildpackage/cmd/danos-buildpackage
```

## Usage

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


