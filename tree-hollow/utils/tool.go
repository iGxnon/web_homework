package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"regexp"
)

// CheckPwdSafe 密码要求: 大于等于6个字符，不能是纯数字，必须有大小写字符和特殊字符，不能有空白字符
// 写的很拉胯
func CheckPwdSafe(password string) bool {

	if len(password) < 6 {
		return false
	}

	// 检测是否是纯数字
	rets := regexp.MustCompile(`\d`)
	alls := rets.FindAllStringSubmatch(password, -1)
	if len(alls) == len(password) {
		return false
	}

	// 检测是否有大写字符
	rets = regexp.MustCompile(`[A-Z]`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) == 0 {
		return false
	}

	// 检测是否有小写字符
	rets = regexp.MustCompile(`[a-z]`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) == 0 {
		return false
	}

	// 检测是否有空白字符
	rets = regexp.MustCompile(`\s`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) != 0 {
		return false
	}

	// 检测是否有特殊字符
	rets = regexp.MustCompile(`\W`)
	alls = rets.FindAllStringSubmatch(password, -1)
	if len(alls) == 0 {
		return false
	}

	return true
}

// PKCS7Padding 把 ciphertext 补成 blockSize 整数倍
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) // 缺多少个字节就补多少个 例：缺5个补5个5
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])   // 拿到最后一个字节(数字)，从而得知之前补了多少个字节
	return origData[:(length - unpadding)] // 获取原本未补齐的 origData
}

// AesEncrypt AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AesDecrypt AES解密
func AesDecrypt(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
