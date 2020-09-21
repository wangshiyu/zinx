package znet

import (
	"github.com/wangshiyu/zinx/utils"
	gorsa "github.com/wangshiyu/zinx/utils/gorsa"
)

type RSA2 struct {
}

//加密
func (r *RSA2) Encryption(data []byte) []byte {
	data, _ = gorsa.PriKeyEncrypt(data, utils.GlobalObject.EncryptionPrivateKey)
	return data
}

//解密
func (r *RSA2) Decrypt(data []byte) []byte {
	data, _ = gorsa.PublicDecrypt(data,utils.GlobalObject.EncryptionPublicKey)
	return data
}

func NewRSA2() *RSA2 {
	return &RSA2{}
}
