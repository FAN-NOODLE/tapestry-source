<?php

namespace App\Modules\Login;

use Illuminate\Support\Facades\Log;

class ScanCodeLogin extends CommonLogin implements LoginInterface
{
    /**
     * 授权码
     *
     * @var string
     */
    private $authCode;

    public function __construct($authCode)
    {
        parent::__construct();

        $this->authCode = $authCode;
    }

    /**
     * @param null|string $token
     * @throws \Exception
     */
    public function init($token = null)
    {
        if (!empty($token)) {
            parent::reduction($token);
        } else {
            // 用auth_code去uc查授权信息
            $authInfo = $this->workEnterpriseSdk->authCodeInfo([
                'auth_code' => $this->authCode
            ]);

            Log::info('用auth_code去uc查授权信息', $authInfo);

            // 判断管理员是否有效
            if (empty($authInfo)) {
                $this->authEffective = false;
            }

            // 说明没有在UC注册企业
            if ($authInfo['enterpriseInfo'] === null) {
                throw new \Exception('系统未查到该企业信息, 请联系企业管理员确认');
            }
            $epId = $authInfo['enterpriseInfo']['epId'];

            // 因为uc没返回corp_id所以调下面的接口去查一下
            $ucEnterpriseConfig = $this->workEnterpriseSdk->configInfo([
                'enterprise_id' => $epId
            ]);

            Log::info('从uc查corp_id', $ucEnterpriseConfig);
            $authInfo['enterpriseInfo']['corp_id'] = $ucEnterpriseConfig['wxqyCorpid'];
            $authInfo['enterpriseInfo']['memUserid'] = $authInfo['adminerInfo']['memUserid'];

            $this->ucEnterprise = $authInfo['enterpriseInfo'];
            $this->enterprise = $this->enterpriseInfo($ucEnterpriseConfig['wxqyCorpid']);
            $this->corpId = $ucEnterpriseConfig['wxqyCorpid'];
            $this->memberId = $authInfo['adminerInfo']['memUserid'] ?? '';
            $this->isExpire = $this->checkExpire();

            $this->admin = $this->adminInfo($this->memberId, $this->enterprise['id'] ?? 0);
        }
    }
}
