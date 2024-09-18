package domain

import (
	"math/rand"
	"regexp"
	"strings"
)

const (
	defaultSuffixLength = 4
	maxSlugLength       = 20
)

// Slugify converts a string to a slug including a random suffix.
func Slugify(input string) string {
	// Convert to lowercase.
	slug := strings.ToLower(input)

	// Replace spaces with hyphens.
	slug = strings.ReplaceAll(slug, " ", "-")

	// Replace German umlauts.
	slug = strings.ReplaceAll(slug, "ä", "ae")
	slug = strings.ReplaceAll(slug, "ö", "oe")
	slug = strings.ReplaceAll(slug, "ü", "ue")
	slug = strings.ReplaceAll(slug, "ß", "ss")

	// Remove any character that is not alphanumeric or hyphen.
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Trim to max length.
	if len(slug) > maxSlugLength {
		slug = slug[:maxSlugLength]
	}

	// Generate random string.
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	randomSuffix := make([]byte, defaultSuffixLength)
	for i := range randomSuffix {
		randomSuffix[i] = charset[rand.Intn(len(charset))]
	}

	// Append random string to slug.
	return slug + "-" + string(randomSuffix)
}
