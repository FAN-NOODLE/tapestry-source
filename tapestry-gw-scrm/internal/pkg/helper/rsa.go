package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func RsaDecode(encryptedBase64 string) string {
	// 生成或加载密钥对
	privateKey := getPrivateKey()
	// 解密
	decrypted, err := decryptMessage(privateKey, encryptedBase64)
	if err != nil {
		fmt.Println("解密失败:", err)
		return ""
	}
	return decrypted
}

// 获取密钥对（如果不存在则生成并保存）
func getKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	// 检查私钥文件是否存在
	if _, err := os.Stat("./internal/pkg/helper/rsa/private.pem"); os.IsNotExist(err) {
		// 生成新密钥对
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Println("密钥生成失败:", err)
			return nil, nil
		}

		// 保存私钥
		savePrivateKeyToFile(privateKey, "./internal/pkg/helper/rsa/private.pem")
		// 保存公钥
		savePublicKeyToFile(&privateKey.PublicKey, "./internal/pkg/helper/rsa/public.pem")

		return privateKey, &privateKey.PublicKey
	}

	// 加载已有密钥
	privateKey, err := loadPrivateKeyFromFile("./internal/pkg/helper/rsa/private.pem")
	if err != nil {
		fmt.Println("私钥加载失败:", err)
		return nil, nil
	}
	return privateKey, &privateKey.PublicKey
}

// 获取私钥
func getPrivateKey() *rsa.PrivateKey {
	// 加载已有密钥
	privateKey, err := loadPrivateKeyFromFile("./internal/pkg/helper/rsa/private.pem")
	if err != nil {
		fmt.Println("私钥加载失败:", err)
		return nil
	}
	return privateKey
}

// 保存私钥到文件
func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err := ioutil.WriteFile(filename, privateKeyPEM, 0600)
	if err != nil {
		fmt.Println("私钥保存失败:", err)
		return
	}
}

// 保存公钥到文件
func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err := ioutil.WriteFile(filename, publicKeyPEM, 0644)
	if err != nil {
		fmt.Println("公钥保存失败:", err)
		return
	}
}

// 从文件加载私钥
func loadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("无法解码PEM块")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// 导出公钥为PEM字符串
func exportPublicKeyToPEM(publicKey *rsa.PublicKey) string {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM)
}

// 解密函数
func decryptMessage(privateKey *rsa.PrivateKey, encryptedBase64 string) (string, error) {
	// Base64解码
	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", nil
	}

	// RSA解密
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
	if err != nil {
		return "", nil
	}

	return string(decrypted), nil
}
