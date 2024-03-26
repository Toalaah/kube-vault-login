package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// JWT is a minimal representation of the core claims required for a JWT.
// Additional claims may be included but are not directly callable.
type JWT struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
	token string
}

// FromFile calls FromString after attempting to read the file at `path`.
func FromFile(path string) (*JWT, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromString(string(b))
}

// FromString reads the contents of `s` into a JWT struct. Any errors
// encountered during decoding are returned to the caller.
func FromString(s string) (*JWT, error) {
	split := strings.Split(s, ".")
	if len(split) != 3 {
		return nil, errors.New("Malformed JWT")
	}
	b, err := base64.RawURLEncoding.DecodeString(split[1])
	if err != nil {
		return nil, err
	}

	jwt := new(JWT)
	if err := json.Unmarshal(b, jwt); err != nil {
		return nil, err
	}

	jwt.token = s
	return jwt, nil
}

// Token returns the "raw" token that was used to populate the JWT.
func (jwt *JWT) Token() string {
	return jwt.token
}

// Expired returns whether this JWT is expired based on its `exp` field.
func (jwt *JWT) Expired() bool {
	t := time.Unix(jwt.Exp, 0)
	return time.Now().After(t)
}

// ExportExecCredential exports the JWT as a kubernetes-compliant
// ExecCredential as described in:
// https://kubernetes.io/docs/reference/config-api/client-authentication.v1/
// The contents of the exec credential are written to `w`.
func (jwt *JWT) ExportExecCredential(w io.Writer) error {
	credential := map[string]interface{}{
		"kind":       "ExecCredential",
		"apiVersion": "client.authentication.k8s.io/v1beta",
		"spec":       struct{}{},
		"status": map[string]string{
			"expirationTimestamp": fmt.Sprint(jwt.Exp),
			"token":               jwt.token,
		},
	}
	b, err := json.Marshal(credential)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
