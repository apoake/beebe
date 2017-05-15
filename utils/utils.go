package utils

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"io"
)

func SHA(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func SaveFile(fileName string, savePath string, format string, reader io.Reader) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(savePath + SHA(fileName) + "." + format, b, 0644); err != nil {
		return err
	}
	return nil
}
