<?php

namespace App\Modules\Login;

use Illuminate\Support\Facades\Log;

class WorkWechatLogin extends CommonLogin implements LoginInterface
{
    protected $token;
    protected $data;

    public function __construct($token, $corpId, $data)
    {
        parent::__construct();

        $this->corpId = $corpId;
        $this->token = $token;
        $this->data = $data;
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
            // 先根据crop_id查UC的企业ID
            $enterpriseInfo = $this->workEnterpriseSdk->detailByCorpId([
                'corp_id' => $this->corpId
            ]);

            Log::info('先根据crop_id查UC的企业ID', $enterpriseInfo);

            // 获取企业号登录信息 找出mem_uid
            $authInfo = $this->workEnterpriseSdk->wxSuiteLoginInfo([
                'token' => $this->token
            ]);

            Log::info('获取企业号登录信息', $authInfo);

            $enterpriseInfo['memUserid'] = $authInfo['adminerInfo']['memUserid'] ?? '';
            $this->ucEnterprise = $enterpriseInfo;
            $this->enterprise = $this->enterpriseInfo($this->corpId);
            //$this->memberId = $this->queryMemberId($authInfo['enterpriseInfo']['epId'], $authInfo['adminerInfo']['eaId']);
            $this->memberId = $authInfo['adminerInfo']['memUserid'] ?? '';
            $this->isExpire = $this->checkExpire();

            $this->admin = $this->adminInfo($this->memberId, $this->enterprise['id'] ?? 0);
        }
    }
}
