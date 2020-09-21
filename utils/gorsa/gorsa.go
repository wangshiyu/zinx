package gorsa

// 公钥加密
func PublicEncrypt(data []byte, publicKey string) ([]byte, error) {
	grsa := RSASecurity{}
	grsa.SetPublicKey(publicKey)

	rsadata, err := grsa.PubKeyENCTYPT(data)
	if err != nil {
		return []byte(""), err
	}
	return rsadata, nil
}

// 私钥加密
func PriKeyEncrypt(data []byte, privateKey string) ([]byte, error) {

	grsa := RSASecurity{}
	grsa.SetPrivateKey(privateKey)

	rsadata, err := grsa.PriKeyENCTYPT(data)
	if err != nil {
		return []byte(""), err
	}
	return rsadata, nil
}

// 公钥解密
func PublicDecrypt(data []byte, publicKey string) ([]byte, error) {
	grsa := RSASecurity{}
	grsa.SetPublicKey(publicKey)

	rsadata, err := grsa.PubKeyDECRYPT(data)
	if err != nil {
		return []byte(""), err
	}
	return rsadata, nil
}

// 私钥解密
func PriKeyDecrypt(data []byte, privateKey string) ([]byte, error) {
	grsa := RSASecurity{}
	grsa.SetPrivateKey(privateKey)

	rsadata, err := grsa.PriKeyDECRYPT(data)
	if err != nil {
		return []byte(""), err
	}
	return rsadata, nil
}
