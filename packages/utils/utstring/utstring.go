package utstring

import (
	"crypto/rand"
	"math/big"
	"project-skbackend/packages/utils/utlogger"
)

func GenerateRandomToken(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var randomString string

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			utlogger.LogError(err)
			return "", err
		}
		randomString += string(charset[n.Int64()])
	}

	// Insert a dash in the middle
	if length > 1 {
		middle := length / 2
		randomString = randomString[:middle] + "-" + randomString[middle:]
	}

	return randomString, nil
}
