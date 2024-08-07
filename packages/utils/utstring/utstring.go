package utstring

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"project-skbackend/packages/utils/utlogger"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomToken(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var (
		randstr string
	)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			utlogger.Error(err)
			return "", err
		}
		randstr += string(charset[n.Int64()])
	}

	// Insert a dash in the middle
	if length > 1 {
		middle := length / 2
		randstr = randstr[:middle] + "-" + randstr[middle:]
	}

	return randstr, nil
}

func PrintJSON(data any) string {
	json, _ := json.MarshalIndent(data, "", "\t")
	utlogger.Info(json)

	return string(json)
}

func AppendName(names ...string) string {
	return strings.Join(names, " ")
}
