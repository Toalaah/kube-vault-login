package main

import (
	"fmt"
	"os"

	"github.com/toalaah/kube-vault-login/cmd"
	"github.com/urfave/cli/v2"
)

var (
	version string = "dev"
	commit  string = ""
	app     *cli.App
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func init() {
	app = &cli.App{
		Name:    "kubectl-vault_login",
		Usage:   "Authenticate to a kubernetes cluster via a vault server's OIDC role endpoint",
		Version: buildVersion(),
		Flags:   []cli.Flag{},
		Commands: cli.Commands{
			cmd.NewGetTokenCmd(),
		},
	}

}

func buildVersion() string {
	v := ""
	v += version
	if commit != "" {
		v += fmt.Sprintf(" (%s)", commit)
	}
	return v
}
