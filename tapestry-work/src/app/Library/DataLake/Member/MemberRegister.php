<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/2 0002
 * Time: 11:06
 */

namespace App\Library\DataLake\Member;



use App\Library\DataLake\Base;

class MemberRegister extends Base
{
    public $errStr = 'register member error';

    public $api = '/api/interface/third/member/register';

    protected $user;
    protected $throttle = true;

    function __construct($user)
    {
        $this->user = $user;
        parent::__construct();
    }

    protected function setQueryParam()
    {

        $this->queryParam = [

            "brand_code"    => $this->config->getBrandCode(),
            "program_code"  => $this->config->getBrandCode(),
            "name"          => $this->user['real_name'],
            "mobile"        => $this->user['mobile'],
            "officePhone"   => "",
            "birthday"      => $this->user['birthday'] ?? "",
            "gender"        => $this->getSex(),
            "email"         => "",
            "chanel_code"   => $this->config->getChannelCode(),
            "store_code"    => $this->config->getStoreCode(),
            "ba_code"       => "", // todo 导购编码,非必填
            "marsMemberNum" => "",
            "RB_Member"     => "",
            "consentStatus" => "0",
            "consentTime"   => $this->getConsentTime(), //同意时间
            "prmtPerson"    => "",
            "address"       => [],
            "media_account" => $this->getMediaAccount()
        ];
    }

    private function getSex()
    {
        switch ((int)$this->user['sex']) {
            case 1:
                return 'M';
            case 2:
                return 'F';
            default:
                return 'U';
        }
    }

    private function getConsentTime()
    {
        return $this->user['register_at'] == 0 ? date("Y-m-d H:i:s"): date("Y-m-d H:i:s", $this->user['register_at'] );
    }

    private function getMediaAccount()
    {
        return [
            [
                "type"           => "3", // 微信
                "accountNo"      => $this->user['open_id'],
                "nickname"       => $this->user['nick_name'],
                "fromSocialCode" => "",
                "remark"         => "",
                "status"         => "0",
                "consentStatus"  => "1",
                "consentTime"    => date('Y-m-d H:i:s', $this->user['created_at']),
                "unionId"        => $this->user['union_id']
            ]
        ];
    }


}