package logging

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateRequestID(prefix string) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", prefix, u.String()), nil
}
