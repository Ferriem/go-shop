package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 16,24,32: aes/128,aes/192,aes/256
var PwdKey = []byte("TESTTESTTESTTEST")

// PKCS7
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	// Repeat returns a new slice consisting of byte(padding).
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS7UnPadding(originData []byte) ([]byte, error) {
	length := len(originData)
	if length == 0 {
		return nil, errors.New("encrypted string error")
	} else {
		unpadding := int(originData[length-1])
		return originData[:(length - unpadding)], nil
	}
}

func AesEncrypt(originData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// PKCS7
	blockSize := block.BlockSize()

	//fill the original data to length of blockSize
	originData = PKCS7Padding(originData, blockSize)

	//CBC mode
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

func AesDeCrypt(crypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(crypted))
	blockMode.CryptBlocks(originData, crypted)

	originData, err = PKCS7UnPadding(originData)
	if err != nil {
		return nil, err
	}

	return originData, nil
}

// base64 encrypt
func EnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(result), nil
}

func DePwdCode(pwd string) ([]byte, error) {
	pwdByte, err := base64.URLEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	return AesDeCrypt(pwdByte, PwdKey)
}
