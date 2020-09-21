package ziface

type IEncryption interface {
	//加密
	Encryption(data []byte)[]byte
	//解密
	Decrypt(data []byte)[]byte
}