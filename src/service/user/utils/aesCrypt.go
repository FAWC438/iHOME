package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

const aes128KeyStr string = "fahsifdaodihhfxp"

// AesCtrCrypt
//
//	@Description: AES-CRT 加密与解密均可采用该函数
//	@param plainText 欲加密/解密的原始数据，byte 类型的切片
//	@return []byte 加密/解密后的数据
//	@return error
func AesCtrCrypt(plainText []byte) ([]byte, error) {
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
