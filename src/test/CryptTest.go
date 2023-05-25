package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// https://www.cnblogs.com/you-men/p/14160439.html

const aes128KeyStr string = "fahsifdaodihhfxp"

// aesCtrCrypt
//
//	@Description: AES-CRT 加密与解密均可采用该函数
//	@param plainText 欲加密/解密的原始数据，byte 类型的切片
//	@return []byte 加密/解密后的数据
//	@return error
func aesCtrCrypt(plainText []byte) ([]byte, error) {
	key := []byte(aes128KeyStr)
	//1. 创建 cipher.Block 接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("5"), block.BlockSize())
	//iv := make([]byte, block.BlockSize())
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	panic(err)
	//}
	//fmt.Println("iv: " + hex.EncodeToString(iv))

	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return dst, nil
}

func hmacSha256(key, data string) string {
	hash := hmac.New(sha256.New, []byte(key)) //创建对应的sha256哈希加密算法
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum([]byte("")))
}

func main() {
	key, data := "me", "kevin"
	fmt.Println(key)
	fmt.Println(data)
	//res := hmacSha256(key, data)
	//fmt.Println(res)

	res, err := aesCtrCrypt([]byte(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("encode: " + string(res))

	res, err = aesCtrCrypt(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("decode: " + string(res))
}
