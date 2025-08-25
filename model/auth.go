package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GenerateSimpleToken(userID uint64) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b) + ":" + fmt.Sprint(userID)
}

func ParseToken(token string) (int64, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid token")
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	return id, err
}
