package utils

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	req "github.com/imroc/req"
	be "github.com/pickjunk/bgo/error"
)

// WxAPI struct
type WxAPI struct {
	Method  string
	URI     string
	Query   req.QueryParam
	Body    map[string]interface{}
	Headers req.Header
}

var wxURL = "https://api.weixin.qq.com"

// Fetch execute a wx api
func (w *WxAPI) Fetch(ctx context.Context, result interface{}) error {
	if os.Getenv("ENV") != "production" {
		req.Debug = true
	}
	defer func() {
		req.Debug = false
	}()

	var res *req.Resp
	var err error
	switch w.Method {
	case "POST":
		res, err = req.Post(wxURL+w.URI, w.Headers, w.Query, req.BodyJSON(w.Body), ctx)
	default:
		res, err = req.Get(wxURL+w.URI, w.Headers, w.Query, ctx)
	}
	if err != nil {
		return err
	}

	code := res.Response().StatusCode
	if !(code >= 200 && code < 300) {
		return fmt.Errorf("http status error: %d", code)
	}

	var e struct {
		Errcode int64
		Errmsg  string
	}
	err = res.ToJSON(&e)
	if err != nil {
		return err
	}
	if e.Errcode > 0 {
		var bErr be.BusinessError
		err = json.Unmarshal([]byte(e.Errmsg), &bErr)
		if err != nil {
			return errors.New(e.Errmsg)
		}
		return &bErr
	}

	if result != nil {
		err = res.ToJSON(result)
		if err != nil {
			return err
		}
	}

	return nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const lettersAndDigits = letters + digits

func randStr(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = lettersAndDigits[rand.Intn(len(lettersAndDigits))]
	}
	return b
}

// pkcs7Pad 补位
// thanks to https://studygolang.com/articles/4752
func pkcs7Pad(data []byte, blockSize int) []byte {
	count := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(count)}, count)
	return append(data, padding...)
}

// pkcs7Unpad 取消补位
// thanks to https://studygolang.com/articles/4752
func pkcs7Unpad(data []byte, blockSize int) []byte {
	dataLen := len(data)
	unpadding := int(data[len(data)-1])
	return data[:(dataLen - unpadding)]
}

// WxEncrypt 第三方平台消息加密
func WxEncrypt(data []byte, platformKey, platformAppid string) ([]byte, error) {
	// base64 decode 密钥
	key := platformKey + "="
	ekey, _ := base64.StdEncoding.DecodeString(key)

	// 把data的长度转化为网络字节序（大端）
	dataLen := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)))

	// 拼装
	var pack bytes.Buffer
	pack.Write(randStr(16))
	pack.Write(dataLen)
	pack.Write(data)
	pack.Write([]byte(platformAppid))

	// 加密前，补位
	target := pkcs7Pad(pack.Bytes(), 32)

	// AES CBC 加密
	block, err := aes.NewCipher(ekey)
	if err != nil {
		return nil, err
	}
	// ekey[:16] 这里为iv，类似于session key
	// 这里很可能是微信算法不规范
	// 一般来说iv应该是随机且要放在加密体里
	cbc := cipher.NewCBCEncrypter(block, ekey[:16])
	cipherData := make([]byte, len(target))
	cbc.CryptBlocks(cipherData, target)

	return cipherData, nil
}

// WxDecrypt 第三方平台消息解密
func WxDecrypt(data []byte, platformKey string) ([]byte, error) {
	// base64 decode 密钥
	key := platformKey + "="
	ekey, _ := base64.StdEncoding.DecodeString(key)

	// AES CBC 解密
	block, err := aes.NewCipher(ekey)
	if err != nil {
		return nil, err
	}
	// ekey[:16] 这里为iv，类似于session key
	// 这里很可能是微信算法不规范
	// 一般来说iv应该是随机且要放在加密体里
	cbc := cipher.NewCBCDecrypter(block, ekey[:16])
	target := make([]byte, len(data))
	cbc.CryptBlocks(target, data)

	// 解密后，取消补位
	target = pkcs7Unpad(target, 32)

	// 提取data
	dataLen := binary.BigEndian.Uint32(target[16:20])
	return target[20 : 20+dataLen], nil
}

// WxSign 第三方平台签名
func WxSign(timestamp, nonce, encrypt, token string) string {
	s := []string{token, timestamp, nonce, encrypt}
	sort.Strings(s)
	str := strings.Join(s, "")
	hash := sha1.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
