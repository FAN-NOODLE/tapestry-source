<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 18:03
 */

namespace App\Http\Controllers\App\V1\PartnerMember;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\Decrypt;
use App\Services\PartnerToken;
use App\Services\Wechat\WxDecrypt;
use App\Traits\CodeTrait;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Work\Member;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\PartnerEnterpriseMap;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;


class RegisterController extends BaseController
{
    protected $isRequireLogin = false;
    const   REGISTER = 1;//去注册
    const   SECOND_REG = 2;//二次注册
    const   REGISTERED = 3;//已注册
    use CodeTrait;

    /**
     * 商城小程序授权类型为：1
     */
    const AUTH_TYPE = 1;
    const SENSITIVE = 1; //脱敏显示
    const PARTNER_SN = 'FX';//合伙人编码前缀
    const USER_TYPE = 2;//合伙人

    const  TIPS = '您已经是%s的企业同事了，不能注册成为合伙人哦';


    public function index(Request $request)
    {
        $rules    = [
            'avatar'           => 'required',//头像
            'partner_name'     => 'required',//合伙人昵称
            'partner_mobile'   => 'required',//手机号
            'mobile_sensitive' => 'required|in:0,1',//手机号是否脱敏显示
            'cache_key'        => 'required',//获得的微信用户的数据信息
        ];
        $messages = [
            'avatar.*'           => '头像信息必须存在',
            'partner_mobile.*'   => '手机号不合法',
            'partner_name.*'     => '合伙人昵称必须存在',
            'mobile_sensitive.*' => '脱敏显示状态必须存在',
            'cache_key.*'        => 'cache_key必须存在',

        ];
        $this->validate($request, $rules, $messages);

        $input = $this->request->all();

        //解参
        $cache = (new Decrypt())->index($input['cache_key']);
        //已注册
        if ($cache['status'] == self::REGISTERED) {
            error(657001);
        }

        //查询导购信息
        $workerInfo = (new StoreWorker())->detail(
            [
                'mem_uid'       => $cache['mem_uid'],
                'enterprise_id' => $cache['enterprise_id'],
                'store_id'      => $cache['store_id'] ?? 0,
                'worker_id'     => $cache['worker_id'],
            ]
        );


        /* //查询合伙人是否是导购，并是否和当前推广导购在一个企业
         $workerInfoByPartnerMobile = (new StoreWorker())->searchStoreWorker(
             [
                 'mobile' => $input['partner_mobile'],
                 'enterprise_id' => $cache['enterprise_id'],
             ]
         );*/
        //导购可以邀请自己
        if ($workerInfo['mobile'] != $input['partner_mobile']) {
            //查询合伙人是否在企业组织里
            // 验证手机号
            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $cache['enterprise_id']]);
            if (empty($enterpriseInfo)) {
                notice(657009);
            }

            $memParams = [
                'qy_domain'      => $enterpriseInfo['data']['domain'],
                'app_identifier' => $enterpriseInfo['data']['identifier'],
                'enterprise_id'  => $cache['enterprise_id'],
                'params'         => ['memMobile' => $input['partner_mobile']]
            ];

            $memberInfo = (new Member())->memberCheck($memParams);
            unset($memberInfo['requestId']);
            $memActive = true;
            if (isset($memberInfo['memActive'])) {
                $memActive = ($memberInfo['memActive'] == 1) ? true : false;
            }
            if (!empty($memberInfo) && $memActive) {
                notice(657003, [$enterpriseInfo['data']['name']]);
            }
        }

        $partner_id = 0;
        //首次注册
        if ($cache['status'] == self::REGISTER) {
            //验证码效验
            if (!isset($input['sms_code'])) {
                error(657002);
            }
            $isSuperUser = in_array($input['partner_mobile'], Constant::WX_LOGIN_SUPER_USER);
            $resultCode  = $this->getCodeCache($input['partner_mobile'], Constant::APP_CLIENT_NAME);
            if ($input['sms_code'] != $resultCode && !$isSuperUser) {
                error(652001);
            }

            //合伙人编码生成 时间戳后6位，加两个随机数
            $randNumber = substr(time(), -6) . rand(0, 9) . rand(0, 9);
            $partner_sn = self::PARTNER_SN . $randNumber;
            //效验
            $partnerInfo = [
                'partner_sn'     => $partner_sn,
                'partner_name'   => $input['partner_name'],
                //'partner_sex'=>$cache['sex'],
                'partner_sex'    => $input['partner_sexex'] ?? 0,
                'partner_mobile' => $input['partner_mobile'],
                'open_id'        => $cache['open_id'],
                //'union_id'=>$cache['union_id'],
                'avatar'         => $input['avatar'],
                'mobile_short'   => substr($input['partner_mobile'], -4)
            ];
            //如果用户微信授权
            if (isset($input['iv']) && $input['iv'] && isset($input['encryptedData']) && $input['encryptedData']) {
                $wx_result                  = (new WxDecrypt($cache['sessionKey'], $this->identifier))->index($input);
                $wx_result                  = json_decode($wx_result, true);
                $partnerInfo['partner_sex'] = $wx_result['gender'];
                //$partnerInfo['union_id'] = $wx_result['unionId'];
            }

            [$appid, $secret] = Decrypt::getEnvMiniProgram($this->identifier);
            $partnerInfo['app_id'] = $appid;
            //添加合伙人基础信息
            try {
                $partnerResult = (new PartnerMember())->insertData($partnerInfo);
            } catch (Error $error) {
                Log::info("添加合伙人基础信息，错误信息：{$error->getMessage()}, 错误code：{$error->getCode()}", $partnerInfo);
                notice(657010);
            }
            $partner_id = $partnerResult['data'];
        }
        //二次注册
        if (isset($input['partner_id']) && $input['partner_id'] && $cache['status'] == self::SECOND_REG) {
            //获得合伙人数据
            $detail     = (new PartnerMember())->detail(['partner_id' => $input['partner_id']]);
            $partner_id = $detail['id'];
        }

        /* //依据union_id 获得商城合伙人user_id
         $userId = (new User())->detail(
             [
                 'union_id'=>$union_id,
                 'enterprise_id'=>$input['enterprise_id']
             ]
         );*/
        $partner_mobile = $input['partner_mobile'];
        if ($input['mobile_sensitive'] == self::SENSITIVE) {
            $partner_mobile = $this->dataDesensitization($input['partner_mobile'], 3, 4);
        }
        $partnerMap = [
            'partner_id'       => $partner_id,
            'enterprise_id'    => $cache['enterprise_id'],
            'worker_id'        => $workerInfo['id'],
            //'user_id'=>isset($userId['data']['id']) ? $userId['data']['id'] : 0,
            'worker_name'      => $workerInfo['name'],
            'worker_uid'       => $workerInfo['mem_uid'],
            'worker_avatar'    => $workerInfo['avatar'],
            'mobile_sensitive' => $input['mobile_sensitive'],//脱敏状态
            'mobile_show'      => $partner_mobile,//根据是否脱敏显示手机号
        ];
        //添加导购合伙人关联表
        try {
            $data = (new PartnerEnterpriseMap())->insertData($partnerMap);
        } catch (Error $error) {
            if ($error->getCode() == 1100032) {
                notice(657011);
            }
            if ($error->getCode() == 1100030) {
                notice(657010);
            }
        }

        $partnerMap['id'] = $data['data'];
        //$partnerData =  array_merge($partnerInfo,$partnerMap);
        //生成token
        $enterpriseParams = [
            'partner_id'    => $partner_id,
            'enterprise_id' => $cache['enterprise_id'],
        ];
        $token            = $this->defaultLogin($enterpriseParams, $input);
        return success(['token' => $token]);
    }

    /**
     * 注册完默认登录获取token
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/25 10:23
     */
    public function defaultLogin($enterpriseParams, $input)
    {
        $loginId     = md5(uniqid() . microtime(true));
        $getToken    = new PartnerToken($input['partner_mobile']);
        $partnerData = (new PartnerMember())->partnerMy($enterpriseParams)['enterprise_info'];
        $partnerData = $getToken->formData($partnerData, $enterpriseParams, $input);
        $adminList   = $partnerData;
        /*foreach ($adminList as $enterprise => $role) {
            $adminList[$enterprise]['guide_roles'] = [
                'role_name' => '分销合伙人',
                'role_mark' => Constant::LOGIN_ROLE_STORE_PARTNER
            ];
        }*/
        foreach ($partnerData as $enterprise => $role) {
            $partnerData[$enterprise]['guide_roles'][] = [
                'role_name' => '分销合伙人',
                'role_mark' => Constant::LOGIN_ROLE_STORE_PARTNER
            ];
        }

        $token = $getToken->partnerToken(self::USER_TYPE, $loginId);
        $getToken->loginCache($adminList, $partnerData, $loginId);
        return $token;
    }

    function dataDesensitization($string, $start = 0, $length = 0, $re = '*')
    {
        if (empty($string)) {
            return false;
        }
        $strarr    = array();
        $mb_strlen = mb_strlen($string);
        while ($mb_strlen) {//循环把字符串变为数组
            $strarr[]  = mb_substr($string, 0, 1, 'utf8');
            $string    = mb_substr($string, 1, $mb_strlen, 'utf8');
            $mb_strlen = mb_strlen($string);
        }
        $strlen = count($strarr);
        $begin  = $start >= 0 ? $start : ($strlen - abs($start));
        $end    = $last = $strlen - 1;
        if ($length > 0) {
            $end = $begin + $length - 1;
        } elseif ($length < 0) {
            $end -= abs($length);
        }
        for ($i = $begin; $i <= $end; $i++) {
            $strarr[$i] = $re;
        }
        if ($begin >= $end || $begin >= $last || $end > $last) return false;
        return implode('', $strarr);
    }


}