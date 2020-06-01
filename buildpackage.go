package buildpackage

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
)

type Builder struct {
	cli             *client.Client
	ctx             context.Context
	version         string
	dockerImageName string
	localImage      bool
	srcDir          string
	destDir         string
	pkgDir          string

	removeContainer bool

	containerID string
}

type MakeBuilderOption func(*Builder)

func WithClient(cli *client.Client) MakeBuilderOption {
	return func(b *Builder) {
		b.cli = cli
	}
}

func WithContext(ctx context.Context) MakeBuilderOption {
	return func(b *Builder) {
		b.ctx = ctx
	}
}

func Version(version string) MakeBuilderOption {
	return func(b *Builder) {
		b.version = version
	}
}
func ImageName(name string) MakeBuilderOption {
	return func(b *Builder) {
		b.dockerImageName = name
	}
}

func LocalImage() MakeBuilderOption {
	return func(b *Builder) {
		b.localImage = true
	}
}

func SourceDirectory(srcDir string) MakeBuilderOption {
	return func(b *Builder) {
		b.srcDir = srcDir
	}
}

func DestinationDirectory(destDir string) MakeBuilderOption {
	return func(b *Builder) {
		b.destDir = destDir
	}
}

func PreferredPackageDirectory(pkgDir string) MakeBuilderOption {
	return func(b *Builder) {
		b.pkgDir = pkgDir
	}
}

func RemoveContainer(remove bool) MakeBuilderOption {
	return func(b *Builder) {
		b.removeContainer = remove
	}
}

func MakeBuilder(opts ...MakeBuilderOption) (*Builder, error) {
	b := new(Builder)
	b.version = "latest"
	b.removeContainer = true
	b.dockerImageName = "jsouthworth/danos-buildpackage"
	for _, opt := range opts {
		opt(b)
	}
	if b.cli == nil {
		cli, err := client.NewEnvClient()
		if err != nil {
			return nil, err
		}
		b.cli = cli
	}
	if b.ctx == nil {
		b.ctx = context.Background()
	}
	if b.srcDir == "" {
		return nil, errors.New("must supply Source Directory")
	}
	if !b.srcDirIsDebian() {
		return nil, errors.New("must be run from the top level of a debian package tree")
	}
	if b.destDir == "" {
		return nil, errors.New("must supply Destination Directory")
	}
	return b, nil
}

func (b *Builder) srcDirIsDebian() bool {
	_, err := os.Stat(b.srcDir + "/debian/control")
	return !os.IsNotExist(err)
}

func (b *Builder) canonicalImageName() string {
	if b.localImage {
		return b.imageName()
	}
	return "registry.hub.docker.com/" + b.imageName()
}

func (b *Builder) imageName() string {
	return b.dockerImageName + ":" + b.version
}

func (b *Builder) pullEnvironment() error {
	if b.localImage {
		return nil
	}
	log.Println("pulling environment", b.canonicalImageName())
	r, err := b.cli.ImagePull(
		b.ctx,
		b.canonicalImageName(),
		types.ImagePullOptions{},
	)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, r)
	return nil
}

func (b *Builder) getBindMounts() []string {
	binds := []string{b.srcDir + ":/mnt/src", b.destDir + ":/mnt/output"}
	if b.pkgDir != "" {
		binds = append(binds, b.pkgDir+":/mnt/pkgs")
	}
	return binds
}

func (b *Builder) createEnvironment() error {
	log.Println("creating environment", b.canonicalImageName())
	createResp, err := b.cli.ContainerCreate(
		b.ctx,
		&container.Config{
			Image:        b.canonicalImageName(),
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{
			Binds: b.getBindMounts(),
		},
		nil,
		"",
	)
	if err != nil {
		return err
	}
	b.containerID = createResp.ID
	log.Println("containerID", b.containerID)
	return nil
}

func (b *Builder) buildPackage() error {
	log.Println("building package", b.containerID)
	out, err := b.cli.ContainerAttach(
		b.ctx,
		b.containerID,
		types.ContainerAttachOptions{
			Stream: true,
			Stdout: true,
			Stderr: true,
		},
	)
	if err != nil {
		return err
	}
	defer out.Close()
	go stdcopy.StdCopy(os.Stdout, os.Stderr, out.Reader)

	err = b.cli.ContainerStart(
		b.ctx,
		b.containerID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	_, err = b.cli.ContainerWait(b.ctx, b.containerID)
	if err != nil {
		return err
	}
	data, err := b.cli.ContainerInspect(b.ctx, b.containerID)
	if err != nil {
		return err
	}
	state := data.State
	if state.ExitCode == 0 {
		return nil
	}
	return errors.New("Build failed: " + state.Error)
}

func (b *Builder) deleteEnvironment() error {
	if b.containerID == "" {
		return nil
	}
	if !b.removeContainer {
		log.Println("preserving environment", b.containerID)
		return nil
	}
	log.Println("deleting environment", b.containerID)
	return b.cli.ContainerRemove(
		b.ctx,
		b.containerID,
		types.ContainerRemoveOptions{},
	)
}

func (b *Builder) Build() (err error) {
	type buildStep func() error

	defer func() {
		e := b.deleteEnvironment()
		if err == nil {
			err = e
		}
	}()

	steps := []buildStep{
		b.pullEnvironment,
		b.createEnvironment,
		b.buildPackage,
	}
	for _, step := range steps {
		err := step()
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) Close() error {
	if b.cli != nil {
		return b.cli.Close()
	}
	return nil
}
