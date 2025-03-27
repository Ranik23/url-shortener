package servicehelpers


import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)


func GenereateShortenedLink(default_link string) (string, error) {
	if default_link == "" {
		return "", errors.New("empty URL not allowed")
	}
	hash := sha256.Sum256([]byte(default_link))
	shortURL := base64.RawURLEncoding.EncodeToString(hash[:])[:8]
	return shortURL, nil
}

