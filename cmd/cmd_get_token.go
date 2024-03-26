package cmd

import (
	"os"

	"github.com/toalaah/kube-vault-login/internal/hash"
	"github.com/toalaah/kube-vault-login/internal/jwt"
	"github.com/toalaah/kube-vault-login/internal/vault"
	"github.com/urfave/cli/v2"
)

var cmdGetToken = &cli.Command{
	Name:  "get-token",
	Usage: "Obtain an OIDC token from the vault server",
	Description: `
Obtain an OIDC token for a specified role from a vault server. The returned
JWT is injected into a kubernetes ExecCredential object and is printed to
stdout.

JWTs returned from the vault server are cached. Subsequent calls of this
subcommand using the same role type will prefer to pull from the cache.
This behaviour can be overwritten.
  `,
	Args: false,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "role",
			Required: true,
			EnvVars: []string{
				"VAULT_OIDC_ROLE",
			},
		},
		&cli.BoolFlag{
			Name: "force-refresh",
			EnvVars: []string{
				"VAULT_OIDC_FORCE_REFRESH",
			},
		},
	},
	Action: func(ctx *cli.Context) error {
		c, err := vault.NewClient()
		if err != nil {
			return err
		}

		role := vault.OIDCRole(ctx.String("role"))
		force := ctx.Bool("force-refresh")

		// If user has not explicitly requested to force-refresh, check if JWT is
		// in cache, if so immediately export it to ExecCredential and return
		if !force {
			cachePath, err := hash.CachePath(role)
			if err != nil {
				return err
			}

			cachedJWT, err := jwt.FromFile(cachePath)
			if err != nil && !os.IsNotExist(err) {
				return err
			}

			if cachedJWT != nil && !cachedJWT.Expired() {
				return cachedJWT.ExportExecCredential(os.Stdout)
			}
		}

		s, err := c.RequestJWT(role)
		if err != nil {
			return err
		}

		newJWT, err := jwt.FromString(s)
		if err != nil {
			return err
		}

		if err := newJWT.ExportExecCredential(os.Stdout); err != nil {
			return err
		}

		return hash.CacheContents(role, []byte(newJWT.Token()))
	},
}

func NewGetTokenCmd() *cli.Command { return cmdGetToken }
