package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	bpkg "jsouthworth.net/go/danos-buildpackage"
)

var srcDir, destDir, pkgDir, version, imageName string
var noClean, localImage bool

func init() {
	flag.StringVar(&srcDir, "src", ".", "source directory")
	flag.StringVar(&destDir, "dest", "..", "destination directory")
	flag.StringVar(&pkgDir, "pkg", "", "preferred package directory")
	flag.StringVar(&version, "version", "latest", "version of danos to build for")
	flag.BoolVar(&noClean, "no-clean", false, "don't delete the container when done")
	flag.StringVar(&imageName, "image-name", "jsouthworth/danos-buildpackage",
		"name of docker image")
	flag.BoolVar(&localImage, "local", false, "is the image only on the local system")
}

func resolvePath(in string) string {
	out, err := filepath.Abs(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return out
}

func main() {
	flag.Parse()
	opts := []bpkg.MakeBuilderOption{
		bpkg.SourceDirectory(resolvePath(srcDir)),
		bpkg.DestinationDirectory(resolvePath(destDir)),
		bpkg.PreferredPackageDirectory(resolvePath(pkgDir)),
		bpkg.RemoveContainer(!noClean),
		bpkg.ImageName(imageName),
		bpkg.Version(version),
	}
	if localImage {
		opts = append(opts, bpkg.LocalImage())
	}
	b, err := bpkg.MakeBuilder(opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}
	err = b.Build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
