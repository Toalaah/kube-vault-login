package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/toalaah/kube-vault-login/cmd"
	"github.com/urfave/cli/v2"
)

var (
	version string
	commit  string
	branch  string

	app *cli.App
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func init() {
	app = &cli.App{
		Name:    "kube-vault-login",
		Usage:   "Authenticate to a kubernetes cluster via a vault server's OIDC role endpoint",
		Version: buildVersionString(),
		Flags:   []cli.Flag{},
		Commands: cli.Commands{
			cmd.NewGetTokenCmd(),
		},
	}

}

func buildVersionString() string {
	dirty := false
	v := version
	dbg, ok := debug.ReadBuildInfo()

	// Set version only if it was not set via ldflags
	if v == "" && ok {
		v = dbg.Main.Version
	}
	// Fallback to unknown default version identifier if ldflags not set or we are in debug context.
	if v == "" || v == "(devel)" {
		v = "dev"
	}

	// Try to read some vcs info from debug build
	if commit == "" && ok {
		for _, setting := range dbg.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value[:7]
			case "vcs.modified":
				if val, err := strconv.ParseBool(setting.Value); err == nil {
					dirty = val
				}
			}
		}
	}

	if dirty {
		v += "-dirty"
	}

	if commit != "" {
		switch branch {
		case "":
			v = fmt.Sprintf("%s (%s)", v, commit)
		default:
			v = fmt.Sprintf("%s (%s %s)", v, commit, branch)
		}
	}

	return v
}
