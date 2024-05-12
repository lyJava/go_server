package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
)

//https://piaohua.github.io/post/golang/20210306-golang-rsa-oaep/

// GenerateRSAKey2File 生成密钥
func GenerateRSAKey2File(bits int, publicKeyName, privateKeyName string) (err error) {
	publicKeyWriter := bytes.NewBuffer([]byte{})
	privateKeyWriter := bytes.NewBuffer([]byte{})

	err = GenerateRSAKey(bits, publicKeyWriter, privateKeyWriter)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("公钥写入", publicKeyWriter.String())
	log.Println("私钥写入", privateKeyWriter.String())

	err = os.WriteFile(publicKeyName, publicKeyWriter.Bytes(), 0444)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = os.WriteFile(privateKeyName, privateKeyWriter.Bytes(), 0400)
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func GenerateRSAKey(bits int, publicKeyWriter, privateKeyWriter io.Writer) error {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)

	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509PrivateKey,
	}
	err = pem.Encode(privateKeyWriter, privateBlock)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// generate public key
	publicKey := &privateKey.PublicKey
	x509PublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509PublicKey,
	}

	err = pem.Encode(publicKeyWriter, publicBlock)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// EncryptOAEP 加密
func EncryptOAEP(publicKey, password []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		err := fmt.Errorf("failed to parse certificate PEM")
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey := pub.(*rsa.PublicKey)

	h := sha256.New() // sha1.New() or md5.New()
	return rsa.EncryptOAEP(h, rand.Reader, rsaPublicKey, password, nil)
}

// DecryptOAEP 解密
func DecryptOAEP(privateKey, cipherdata []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		err := fmt.Errorf("failed to parse certificate PEM")
		return nil, err
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) // ASN.1 PKCS#1 DER encoded form.
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	return rsa.DecryptOAEP(h, rand.Reader, priv, cipherdata, nil)
}
