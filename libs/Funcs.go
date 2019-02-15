package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"reflect"
	"regexp"
)

//对密文删除填充
func unPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

//解密
func WxDecryptData(encryptedData string, iv string, key string) []byte {
	aesKey, err := base64.StdEncoding.DecodeString(key)
	if err == nil {
		aesIV, err := base64.StdEncoding.DecodeString(iv)
		if err == nil {
			aesCipher, err := base64.StdEncoding.DecodeString(encryptedData)
			if err == nil {
				block, err := aes.NewCipher(aesKey)
				if err == nil {
					blockMode := cipher.NewCBCDecrypter(block, aesIV)
					plainText := make([]byte, len(aesCipher))
					blockMode.CryptBlocks(plainText, aesCipher)
					plainText = unPadding(plainText)
					if err == nil {
						return plainText
					}
				}
			}
		}
	}
	return []byte{}
}

func CleanScriptTags(content string) string {
	re := regexp.MustCompile(`@<script[^>]*?>.*?</script>@si`)
	content = re.ReplaceAllString(content, "")
	re = regexp.MustCompile(`@<style[^>]*?>.*?</style>@siU`)
	content = re.ReplaceAllString(content, "")
	re = regexp.MustCompile(`@<![\s\S]*?--[ \t\n\r]*>@`)
	content = re.ReplaceAllString(content, "")
	return content
}

func StripTags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)
	return re.ReplaceAllString(content, "")
}
