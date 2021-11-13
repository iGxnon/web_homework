package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// 有一些功能懒得做了（
// 加密方法用了RSA非对称加密，客户端用公钥加密后传输给服务端，服务端保存加密后的密码到数据池，并持有私钥可以获取密文
// 用非对称加密，有效防止密码被窃取，数据池内仅保存加密后的密码，防止数据池被渗透获取到内容

// Golang的设计方式和Java有很大的不同...

func main() {
	pool := &DataPool{
		path: "/Users/igxnon/个人项目/Golang/web_homework/homework5/level03/users_server.data",
	}
	pool.ReloadAll()
	var loginInfo = LoginPage{}
	loginInfo.Login_1()
	loginInfo.Login_2(*pool)
	fmt.Println(loginInfo)
}

// 生成密钥对

func GenerateRSAKeyPair(bits int) (public, private string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	defer func() {
		if re := recover(); re != nil {
			fmt.Println(re)
		}
	}()
	if err != nil {
		panic(err)
	}
	// x509 对公钥和私钥序列化成ASN.1 的 DER编码字符串
	private = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
	public = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&(privateKey.PublicKey)))
	return
}

// 公钥加密

func RSA_Encrypt(text []byte, public string) (cipherText []byte, err error) {
	decodePublic, err := base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.ParsePKCS1PublicKey(decodePublic)
	if err != nil {
		return nil, err
	}
	cipherText, err1 := rsa.EncryptPKCS1v15(rand.Reader, publicKey, text)
	if err1 != nil {
		return nil, err
	}
	return cipherText, nil
}

//私钥解密

func RSA_Decrypt(cipherText []byte, private string) (text []byte, err error) {
	decodePrivate, err := base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(decodePrivate)
	if err != nil {
		return nil, err
	}
	plainText, err1 := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err1 != nil {
		return nil, err1
	}
	return plainText, nil
}
