package hash

import (
	"crypto/sha1"
	"encoding/hex"
	// "crypto/md5"
)

// GetHash gets the sha1 hash of a string and returns it
func GetHash (s string) (hash string, err error) {
	h := sha1.New()


	_, err = h.Write([]byte(s))
	if (err != nil) {
		return "", err
	}
	
	bs := h.Sum(nil)
	hash = hex.EncodeToString(bs)
	return hash, nil
}