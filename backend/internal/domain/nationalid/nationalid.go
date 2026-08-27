package nationalid

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const Length = 13

func IsValid(id string) bool {
	if len(id) != Length {
		return false
	}
	sum := 0
	for i := 0; i < Length-1; i++ {
		c := id[i]
		if c < '0' || c > '9' {
			return false
		}
		sum += int(c-'0') * (Length - i)
	}
	last := id[Length-1]
	if last < '0' || last > '9' {
		return false
	}
	return (11-sum%11)%10 == int(last-'0')
}

func Hash(pepper, id string) string {
	m := hmac.New(sha256.New, []byte(pepper))
	m.Write([]byte(id))
	return hex.EncodeToString(m.Sum(nil))
}

func Mask(id string) string {
	if len(id) != Length {
		return ""
	}
	return id[0:1] + "-" + id[1:5] + "-xxxxx-xx-" + id[12:13]
}
