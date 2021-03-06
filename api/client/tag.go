package client

import (
	"errors"

	Cli "github.com/hyperhq/hypercli/cli"
	flag "github.com/hyperhq/hypercli/pkg/mflag"
	"github.com/hyperhq/hypercli/reference"
	"github.com/docker/engine-api/types"
)

// CmdTag tags an image into a repository.
//
// Usage: docker tag [OPTIONS] IMAGE[:TAG] [REGISTRYHOST/][USERNAME/]NAME[:TAG]
func (cli *DockerCli) Tag(args ...string) error {
	cmd := Cli.Subcmd("tag", []string{"IMAGE[:TAG] [REGISTRYHOST/][USERNAME/]NAME[:TAG]"}, Cli.DockerCommands["tag"].Description, true)
	force := cmd.Bool([]string{"#f", "#-force"}, false, "Force the tagging even if there's a conflict")
	cmd.Require(flag.Exact, 2)

	cmd.ParseFlags(args, true)

	ref, err := reference.ParseNamed(cmd.Arg(1))
	if err != nil {
		return err
	}

	if _, isCanonical := ref.(reference.Canonical); isCanonical {
		return errors.New("refusing to create a tag with a digest reference")
	}

	var tag string
	if tagged, isTagged := ref.(reference.NamedTagged); isTagged {
		tag = tagged.Tag()
	}

	options := types.ImageTagOptions{
		ImageID:        cmd.Arg(0),
		RepositoryName: ref.Name(),
		Tag:            tag,
		Force:          *force,
	}

	return cli.client.ImageTag(options)
}
