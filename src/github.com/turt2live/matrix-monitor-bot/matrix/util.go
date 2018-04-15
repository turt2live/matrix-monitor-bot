package matrix

import (
	"strings"
	"github.com/pkg/errors"
)

func ExtractUserHomeserver(userId string) (string, error) {
	if string(userId[0]) != "@" {
		return "", errors.New("User ID does not start with @")
	}

	idx := strings.Index(userId, ":")
	if idx <= 1 {
		return "", errors.New("No localpart for user")
	}

	return userId[idx+1:], nil
}
