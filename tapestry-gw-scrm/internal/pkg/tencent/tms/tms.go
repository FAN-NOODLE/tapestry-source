package tms

import (
	"encoding/base64"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20201229"
	"gw-scrm/internal/config"
)

// Verification 验证文本中是否存在违规汉字
func Verification(content string) ([]string, bool, error) {
	credential := common.NewCredential(
		config.CfgData.Tencent.Tms.SecretId,
		config.CfgData.Tencent.Tms.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"
	client, _ := tms.NewClient(credential, config.CfgData.Tencent.Tms.Region, cpf)
	request := tms.NewTextModerationRequest()
	base64Content := base64.StdEncoding.EncodeToString([]byte(content))
	request.Content = common.StringPtr(base64Content)

	response, err := client.TextModeration(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok || err != nil {
		return []string{}, false, err
	}
	if *response.Response.Suggestion != "Pass" {
		return common.StringValues(response.Response.Keywords), false, err
	}
	return []string{}, true, nil
}
