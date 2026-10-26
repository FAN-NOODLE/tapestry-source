<?php

namespace App\Sdk\Active\CRM\CoachCore;

use App\Exceptions\Error;
use App\Sdk\Active\CRM\CoachCore\Auth\TokenCache;
use App\Sdk\Active\CRM\CoachCore\Welfare\WelfareGive;
use App\Sdk\Active\CRM\CoachCore\Welfare\WelfareHaved;

class WelfareRequestService
{
    public $enterpriseId;

    /**
     * @var Request
     */
    public $request;
    public $requestCore;

    public function __construct()
    {
    }

    public function setOption(array $user): WelfareRequestService
    {
//        if (isset($user['mobile'])) {
//            $this->mobile = $user['mobile'];
//        }
//        if (isset($user['union_id'])) {
//            $this->unionId = $user['union_id'];
//        }
        $this->enterpriseId = $user['enterprise_id'];
        $this->requestCore  = new RequestCore($this->enterpriseId);

        return $this;
    }

    /**
     * 获取token
     *
     * @return mixed|string|null
     */
    public function getToken($enterpriseId)
    {
        return (new TokenCache($enterpriseId))->getAccessToken();
    }


    public function actionWelfareGive($campaignSID, $unionID, $extUserid): array
    {
        $request  = new WelfareGive($campaignSID, $unionID, $extUserid);
        $response = $this->requestCore->requestApiCore($request);
        if ($response && $response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

    public function actionWelfareHaved($campaignSID, $unionID): array
    {
        $request  = new WelfareHaved($campaignSID, $unionID);
        $response = $this->requestCore->requestApiCore($request);
        if ($response && $response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

}