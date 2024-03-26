package vault

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	TokenEndpoint = "identity/oidc/token"
)

type OIDCRole string

func (v OIDCRole) Hash() (string, error) {
	h := md5.New()
	h.Write([]byte(os.Getenv("VAULT_TOKEN")))
	h.Write([]byte(v))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

type Client struct {
	*api.Client
}

// NewClient returns a new vault client. Address and token initialization are
// handled internally. Any errors encountered during initialization (for
// instance due to lacking environment variables) are returned to the caller.
func NewClient() (*Client, error) {
	c, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}
	// Try to read from ~/.vault-token if env var is not supplied
	if os.Getenv("VAULT_TOKEN") == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("VAULT_TOKEN unset and/or failed to read token from ~/.vault-token: %s", err)
		}
		token, err := os.ReadFile(path.Join(homeDir, ".vault-token"))
		if err != nil {
			return nil, fmt.Errorf("VAULT_TOKEN unset and/or failed to read token from ~/.vault-token: %s", err)
		}
		c.SetToken(strings.TrimSuffix(string(token), "\n"))
		os.Setenv("VAULT_TOKEN", c.Token())
	}
	return &Client{c}, nil
}

func (c *Client) RequestJWT(r OIDCRole) (string, error) {
	endpoint := fmt.Sprintf("%s/%s", TokenEndpoint, r)

	res, err := c.Logical().Read(endpoint)
	if err != nil {
		return "", err
	}

	s, ok := res.Data["token"].(string)
	if !ok {
		return "", fmt.Errorf("Role endpoint '%s' did not return field 'token'", endpoint)
	}

	return s, nil
}
