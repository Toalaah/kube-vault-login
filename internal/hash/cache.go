package hash

import (
	"os"
	"path"

	"github.com/adrg/xdg"
)

var (
	cacheFolder = path.Join(xdg.CacheHome, "kube-vault-login")
)

func CachePath(h Hasher) (string, error) {
	hash, err := h.Hash()
	if err != nil {
		return "", err
	}
	return path.Join(cacheFolder, hash), nil
}

func CachePathExists(h Hasher) (string, error) {
	hash, err := h.Hash()
	if err != nil {
		return "", err
	}

	cachePath := path.Join(cacheFolder, hash)
	_, err = os.Stat(cachePath)

	return cachePath, err
}

func CacheContents(h Hasher, contents []byte) error {
	if _, err := os.Stat(cacheFolder); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(cacheFolder, 0700); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	hash, err := h.Hash()
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(cacheFolder, hash), contents, 0600)
}
