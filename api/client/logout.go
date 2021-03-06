package client

import (
	"fmt"

	Cli "github.com/hyperhq/hypercli/cli"
	flag "github.com/hyperhq/hypercli/pkg/mflag"
)

// CmdLogout logs a user out from a Docker registry.
//
// If no server is specified, the user will be logged out from the registry's index server.
//
// Usage: docker logout [SERVER]
func (cli *DockerCli) CmdLogout(args ...string) error {
	cmd := Cli.Subcmd("logout", []string{"[SERVER]"}, Cli.DockerCommands["logout"].Description+".\nIf no server is specified, the default is defined by the daemon.", true)
	cmd.Require(flag.Max, 1)

	cmd.ParseFlags(args, true)

	var serverAddress string
	if len(cmd.Args()) > 0 {
		serverAddress = cmd.Arg(0)
	} else {
		serverAddress = cli.electAuthServer()
	}

	if _, ok := cli.configFile.AuthConfigs[serverAddress]; !ok {
		fmt.Fprintf(cli.out, "Not logged in to %s\n", serverAddress)
		return nil
	}

	fmt.Fprintf(cli.out, "Remove login credentials for %s\n", serverAddress)
	delete(cli.configFile.AuthConfigs, serverAddress)
	if err := cli.configFile.Save(); err != nil {
		return fmt.Errorf("Failed to save docker config: %v", err)
	}

	return nil
}
