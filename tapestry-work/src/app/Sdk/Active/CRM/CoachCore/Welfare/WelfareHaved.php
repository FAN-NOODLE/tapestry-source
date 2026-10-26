<?php

namespace App\Sdk\Active\CRM\CoachCore\Welfare;

use App\Sdk\Active\CRM\CoachCore\Base;
use App\Sdk\Active\CRM\CoachCore\Config;
use App\Sdk\Active\CRM\CryptDes;

/**
 * 是否发放过优惠券
 */
class WelfareHaved extends Base
{
    public $errStr = 'WelfareHaved error';
    public $api = '/Welfare/Haved';


    /**
     * CampaignSID  string  优惠券活动sid ,FC749E29-292A-40D8-B906-247943D1DF6D
     * UnionID  string  外部联系人unionid
     */

    private $CampaignSID;
    private $UnionID;

    public function __construct(string $campaignSID, string $unionID)
    {
        parent::__construct();
        $this->CampaignSID = $campaignSID;
        $this->UnionID     = $unionID;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $queryParam = [
            'CampaignSID' => $this->CampaignSID,
            'UnionID'     => $this->UnionID,
        ];

        $this->queryParam = json_encode($queryParam, JSON_UNESCAPED_UNICODE);

        // DES ECB 加解密
        $des              = new CryptDes(Config::key, 'DES-ECB', CryptDes::OUTPUT_BASE64);
        $this->queryParam = $des->encrypt($this->queryParam);

    }
}