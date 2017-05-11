package utils

import (
	"crypto/sha1"
	"fmt"
)

func SHA(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
