package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Sign sha1 签名.
func Sha1Sign(strs ...string) (signature string) {
	sort.Strings(strs)
	s := strings.Join(strs, "")
	hashsum := sha1.Sum([]byte(s))
	return hex.EncodeToString(hashsum[:])
}

// 微信消息加解密
// http://mp.weixin.qq.com/wiki/6/90f7259c0d0739bbb41d9f4c4c8e59a2.html
// 参考 https://github.com/heroicyang/wechat-crypter

// PKCS7Decode 方法用于删除解密后明文的补位字符
func PKCS7Decode(text []byte) []byte {
	pad := int(text[len(text)-1])
	if pad < 1 || pad > 32 {
		pad = 0
	}
	return text[:len(text)-pad]
}

// PKCS7Encode 方法用于对需要加密的明文进行填充补位
func PKCS7Encode(text []byte) []byte {
	const BlockSize = 32
	amountToPad := BlockSize - len(text)%BlockSize
	for i := 0; i < amountToPad; i++ {
		text = append(text, byte(amountToPad))
	}
	return text
}

// Decrypt 方法用于对密文进行解密
// 返回解密后的消息，错误信息
func (ctx *Context) Decrypt(text string) ([]byte, error) {
	var (
		msgDecrypt []byte
		id         string
		msgLen     int32
	)

	deciphered, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(ctx.nowAESKey)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(c, ctx.nowAESKey[:16])
	cbc.CryptBlocks(deciphered, deciphered)

	decoded := PKCS7Decode(deciphered)

	buf := bytes.NewBuffer(decoded[16:20])

	binary.Read(buf, binary.BigEndian, &msgLen)
	msgDecrypt = decoded[20 : 20+msgLen]
	fmt.Println(string(msgDecrypt))
	id = string(decoded[20+msgLen:])
	if id != ctx.appid {
		err := fmt.Errorf("appid error %s", id)
		return nil, err
	}

	return msgDecrypt, nil
}

// Encrypt 方法用于对明文进行加密
func (ctx *Context) Encrypt(message []byte) (string, error) {
	//message := []byte(text)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(message)))
	if err != nil {
		return "", err
	}

	msgLen := buf.Bytes()

	randBytes := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, randBytes)
	if err != nil {
		return "", err
	}

	messageBytes := bytes.Join([][]byte{randBytes, msgLen, message, []byte(ctx.appid)}, nil)

	encoded := PKCS7Encode(messageBytes)

	c, err := aes.NewCipher(ctx.AESKey)
	if err != nil {
		return "", err
	}

	cbc := cipher.NewCBCEncrypter(c, ctx.AESKey[:16])
	cbc.CryptBlocks(encoded, encoded)

	return base64.StdEncoding.EncodeToString(encoded), nil
}
