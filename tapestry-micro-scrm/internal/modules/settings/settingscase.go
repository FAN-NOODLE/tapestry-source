package settings

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/uid"
)

// Service 应用服务
type Service struct {
	repo Repository
}

// SettingsInfo 获取会话存档配置
func (s Service) SettingsInfo(c ctx.Context, params *SettingsInfoRequest) (*ChatSettings, error) {
	resDetail, err := s.repo.GetSettingsInfo(c, *params)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	if resDetail == nil {
		// 生成共钥私钥
		keyInfo, keyerr := GenRsaKey(2048)
		if errorIgnore.Ignore(keyerr) != nil {
			return nil, keyerr
		}
		settingsInfo := ChatSettings{
			Id:                uid.ID(ctx.Wrap(c)),
			MessagePublicKey:  keyInfo.PublicKey,
			MessagePrivateKey: keyInfo.PrivateKey,
			EnterpriseId:      params.EnterpriseID,
		}
		_, errAdd := s.repo.Add(c, settingsInfo)
		if errorIgnore.Ignore(errAdd) != nil {
			return nil, errAdd
		}
	}
	return resDetail, nil
}

// 生成密钥对
func GenRsaKey(bits int) (*GenRsaKeyInfo, error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	keyData := GenRsaKeyInfo{
		PublicKey:  string(pem.EncodeToMemory(publicBlock)),
		PrivateKey: string(pem.EncodeToMemory(priBlock)),
	}
	keyData.PublicKey = strings.TrimSuffix(keyData.PublicKey, "\n")
	keyData.PrivateKey = strings.TrimSuffix(keyData.PrivateKey, "\n")
	return &keyData, nil
}

// ConfigAdd 添加超时提醒设置
func (s Service) SettingsAdd(c ctx.Context, data ChatSettings) error {
	_, err := s.repo.Add(c, data)
	if err != nil {
		return err
	}
	return nil
}

// ConfigAdd 编辑超时提醒设置
func (s Service) SettingsSave(c ctx.Context, data *ChatSettings) (*ChatSettings, error) {
	_, err := s.repo.Update(c, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// EnterpriseList 企业列表
func (s Service) EnterpriseList(c ctx.Context) []ChatSettings {
	var list []ChatSettings
	list, err := s.repo.EnterpriseList(c)
	if err != nil {
		return list
	}
	return list
}
