package mai

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

func DecHashToRaw(raw string) string {
	formatData := CBCDecrypt(raw, os.Getenv("hashkey"))
	return formatData
}

func CBCDecrypt(ciphertext string, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}
	ciphercode, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return ""
	}
	iv := ciphercode[:aes.BlockSize]
	ciphercode = ciphercode[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphercode, ciphercode)
	plaintext := string(ciphercode)
	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
}

func RawJsonParse(raw string) (qq int64, Session string) {
	getQQ := gjson.Get(raw, "qq").String()
	strToInt, err := strconv.ParseInt(getQQ, 10, 64)
	if err != nil {
		return 0, ""
	}
	getSession := gjson.Get(raw, "session").String()
	return strToInt, getSession
}
