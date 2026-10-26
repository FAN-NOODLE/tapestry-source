<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/19
 * Time: 4:44 下午
 */

namespace App\Services;


use App\Common\Constant;
use Carbon\Carbon;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;

class PartnerToken
{

    protected  $loginCacheKey;//cache_key
    protected  $expiredTime; //过期时间
    protected  $mobile;
    function __construct($mobile)
    {
        $this->expiredTime = Constant::LOGIN_TOKEN_EXPIRED + time();
        $this->loginCacheKey = "wx_login_{$mobile}";
        $this->mobile = $mobile;

    }
    /**
     * 合伙人token
     *
     * @param $mobile)
     * @param $userType
     * @param $partnerLoginData
     * @return string
     */
    public  function partnerToken($userType,$loginId)
    {
        // 合伙人token
        $partnerTokenData = [
            'login_cache_key'         => $this->loginCacheKey,
            'mobile'                  => $this->mobile,
            'loginId'                 => $loginId,
            'expiredTime'             => $this->expiredTime,
            'user_type'               => $userType, //1导购 2合伙人，3两个角色
        ];
        $partnerTokenStr  = Crypt::encrypt(json_encode($partnerTokenData));

        return $partnerTokenStr;
    }

    /**
     * @param $loginCacheKey
     * @param $loginId
     * @param $expiredTime
     * @param $adminList
     * @param $partnerLoginDataNew
     */
    public  function loginCache($adminList,$partnerLoginDataNew,$loginId)
    {
        // 登录ID 存入缓存，以实现同一端账号互顶
        Cache::put($this->loginCacheKey, [
            'login_cache_key'         => $this->loginCacheKey,
            'loginId'                 => $loginId,
            'expiredTime'             => $this->expiredTime,
            'enterprise'              => $adminList,
            'partner_enterprise_list' => array_values($partnerLoginDataNew)
        ], Carbon::createFromTimestamp($this->expiredTime));
        return true;
    }



    /**
     * @param $data
     * @return mixed
     */
    public  function formData($data,$enterpriseParams,$input)
    {
        $partner_id = $enterpriseParams['partner_id'];
        $account_name = $input['partner_name'] ?? '';
        $dataNew = [];
        foreach ($data as $key => &$val) {
            $dataNew[$key]['enterprise_id']          = $val['id'];
            $dataNew[$key]['hash']                   = $val['hash'];
            $dataNew[$key]['version']                = $val['version'];
            $dataNew[$key]['version_tag']            = $val['version_tag'];
            $dataNew[$key]['name']                   = $val['name'];
            $dataNew[$key]['contact_name']           = $val['contact'];
            $dataNew[$key]['contact_mobile']         = $val['mobile'];
            $dataNew[$key]['address']                = $val['address'];
            $dataNew[$key]['crs_status']             = $val['crs_status'];
            $dataNew[$key]['store_status']           = $val['store_status'];
            $dataNew[$key]['contact_status']         = $val['contact_status'];
            $dataNew[$key]['domain']                 = $val['domain'];
            $dataNew[$key]['corp_id']                = $val['corp_id'];
            $dataNew[$key]['account_start_at']       = $val['account_start_at'];
            $dataNew[$key]['account_expire_at']      = $val['account_expire_at'];
            $dataNew[$key]['is_split_flow']          = $val['is_split_flow'];
            $dataNew[$key]['admin_id']               = '';
            $dataNew[$key]['account']                = '+86-' . $this->mobile;
            $dataNew[$key]['partner_status']         = $val['account_status'] ?? 1;
            $dataNew[$key]['account_name']           = $account_name ?? '';
            $dataNew[$key]['account_global_roaming'] = '+86';
            $dataNew[$key]['account_mobile']         = $this->mobile;
            $dataNew[$key]['account_email']          = '';
            $dataNew[$key]['role']                   = 0;
            $dataNew[$key]['is_super']               = 0;
            $dataNew[$key]['range_type']             = 0;
            $dataNew[$key]['admin_available_system'] = '';
            $dataNew[$key]['shop_name']              = $val['name'] ?? '';
            $dataNew[$key]['identifier']             = 1;
            $dataNew[$key]['logo']                   = $val['logo'] ?? '';
            $dataNew[$key]['partner_id']             = $partner_id;
        }

        return $dataNew;
    }


}