package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/keys-pub/keys/encoding"
)

func RandomShortCode() string {
	// Generate Random text
	str := rand.Text()

	// Compute SHA-256 hash and return as Base 62
	sum := sha256.Sum256([]byte(str))
	return encoding.EncodeBase62(sum[:])
}

func GenerateTruncatedShortCode(destinationURL string, length int) (string, error) {
	code, err := GenerateShortCode(destinationURL)
	if err != nil {
		return "", err
	}

	// Truncate to 8 characters
	return code[:length], nil
}

func GenerateShortCode(destinationURL string) (string, error) {
	// Canonicalise URL
	u, err := CanonicaliseURL(destinationURL)
	if err != nil {
		return "", err
	}

	// Compute SHA-256 Hash of URL
	sum := sha256.Sum256([]byte(u))

	// Base 62 Encode
	str := encoding.EncodeBase62(sum[:])
	return str, nil
}

// CanonicaliseURL returns a canonical URL or an error if the rawURL is invalid. Normalisation steps:
//
//  1. Forces lowercase scheme and host
//  2. Cleans the path
//  3. Removes trailing slashes
//  4. Standardises query parameters and sorts alphabetically.
//  5. Removes default ports
//  6. Removes fragments
func CanonicaliseURL(rawURL string) (string, error) {
	// Initial URL parse & Validation
	u, err := ValidateAndParseURL(rawURL)
	if err != nil {
		return "", err
	}

	// Normalise the scheme and host to lowercase
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// Clean the path while preserving an empty path for host-only URLs.
	if u.Path != "" {
		u.Path = path.Clean(u.Path)
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	// Standardize query parameters alphabetically
	u.RawQuery = u.Query().Encode()
	u.Fragment = ""

	// Remove default ports
	if (u.Scheme == "http" && strings.HasSuffix(u.Host, ":80")) ||
		(u.Scheme == "https" && strings.HasSuffix(u.Host, ":443")) {
		u.Host = strings.Split(u.Host, ":")[0]
	}

	return u.String(), nil
}

func ValidateAndParseURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "https" {
		return nil, fmt.Errorf("URL must use HTTPS. URL: %s", rawURL)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("URL must contain a host. URL: %s", rawURL)
	}

	return u, nil
}

func ValidateURL(rawURL string) error {
	_, err := ValidateAndParseURL(rawURL)
	return err
}
