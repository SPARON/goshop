package utils

//加密解密工具

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
)

//AesEncrypt Aes加密
func AesEncrypt(origData []byte, key ...[]byte) ([]byte, error) {
	srckey := Conf.SecretKey
	if len(key) > 0 && key[0] != nil {
		srckey = key[0]
	}
	srckey = MD5(srckey)
	block, err := aes.NewCipher(srckey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, srckey[:blockSize])
	cypted := make([]byte, len(origData))
	blockMode.CryptBlocks(cypted, origData)
	return []byte(hex.EncodeToString(cypted)), nil
}

//MD5 MD5加密
func MD5(str []byte) []byte {
	md5code := md5.New()
	md5code.Write(str)
	return md5code.Sum(nil)[:]
}

//AesDecrypt Aes解密
func AesDecrypt(cypted []byte, key ...[]byte) ([]byte, error) {
	cypted, err := hex.DecodeString(string(cypted))
	if err != nil {
		return nil, err
	}
	srckey := Conf.SecretKey
	if len(key) > 0 && key[0] != nil {
		srckey = key[0]
	}
	srckey = MD5(srckey)
	block, err := aes.NewCipher(srckey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, srckey[:blockSize])
	origData := make([]byte, len(cypted))
	blockMode.CryptBlocks(origData, cypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//PKCS5Padding 密钥填充方式
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding 密钥填充方式返填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
