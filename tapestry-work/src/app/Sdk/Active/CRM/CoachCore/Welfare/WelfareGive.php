<?php

namespace App\Sdk\Active\CRM\CoachCore\Welfare;

use App\Sdk\Active\CRM\CoachCore\Base;
use App\Sdk\Active\CRM\CoachCore\Config;
use App\Sdk\Active\CRM\CryptDes;

/**
 * 发放优惠券
 */
class WelfareGive extends Base
{
    public $errStr = 'WelfareGive error';
    public $api = '/Welfare/Give';


    /**
     * CampaignSID  string  优惠券活动sid ,FC749E29-292A-40D8-B906-247943D1DF6D
     * UnionID  string  外部联系人unionid
     * ExtUserid  string 外部联系人userid
     */

    private $CampaignSID;
    private $UnionID;
    private $ExtUserid;

    public function __construct(string $campaignSID, string $unionID, string $extUserid)
    {
        parent::__construct();
        $this->CampaignSID = $campaignSID;
        $this->UnionID     = $unionID;
        $this->ExtUserid   = $extUserid;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {

        $queryParam = [
            'CampaignSID' => $this->CampaignSID,
            'UnionID'     => $this->UnionID,
            'ExtUserid'   => $this->ExtUserid,
        ];

        $this->queryParam = json_encode($queryParam, JSON_UNESCAPED_UNICODE);

//        $this->queryParam = '{"OrderNo":"1111","Url":"123.com","Content":"快点付钱哦!要过期了"}';

        // DES ECB 加解密
        $des = new CryptDes(Config::key, 'DES-ECB', CryptDes::OUTPUT_BASE64);
        $this->queryParam = $des->encrypt($this->queryParam);

    }
}