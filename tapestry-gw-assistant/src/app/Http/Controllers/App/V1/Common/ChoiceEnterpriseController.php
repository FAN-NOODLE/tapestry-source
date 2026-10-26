<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-18
 * Time: 17:55
 */

namespace App\Http\Controllers\App\V1\Common;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\MicroScrm\GoMicroScrm;
use Carbon\Carbon;
use CrsSdk\MicroAuthority\Account;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Validation\ValidationException;

class ChoiceEnterpriseController extends Controller
{

    // 企业二级分销状态
    const PARTNER_STATUS_OFF = 0;

    protected $isRequireLogin = false;

    /**
     * 选择企业
     *
     * @param Request $request
     * @param Account $account
     *
     * @return JsonResponse|void
     * @throws ValidationException
     * @throws Exception
     */
    public function __invoke(Request $request, Account $account)
    {
        $requestData = $this->validate($request, [
            'enterprise_id'   => 'required',
            'login_cache_key' => 'required'
        ], [
            'enterprise_id.required'   => '企业ID不能为空',
            'login_cache_key.required' => '缓存Key不能为空'
        ]);
        // 登录的Token 是否失效
        $cacheData = Cache::get($requestData['login_cache_key']);
        if ($cacheData['expiredTime'] < time()) {
            error(652004);
        }

        $loginUserType = $this->loginUserType();

        // 选择的企业的信息
        $memberInfoList = $cacheData['enterprise'] ?? [];
        if ($loginUserType == self::USERTYPE_PARTNER) {
            $memberInfoList = $cacheData['partner_enterprise_list'] ?? [];
        }

        $memberInfoList  = array_column($memberInfoList, null, 'enterprise_id');
        $loginEnterprise = $memberInfoList[$requestData['enterprise_id']] ?? [];

//        $loginRole = '';
//        if (!empty($loginEnterprise['guide_roles'])) {
//            foreach ($loginEnterprise['guide_roles'] as $role) {
//                if (in_array($role['role_mark'], (array)Constant::LOGIN_ROLE_LIST)) {
//                    $loginRole = $role['role_mark'];
//                    break;
//                }
//            }
//        }
//
//        if (empty($loginRole)) {
//           // notice('login_manager_no_authority');
//        }

//        //合伙人字段映射
//        if (isset($loginEnterprise['partner_mobile'])) {
//            $loginEnterprise['account_mobile'] = $loginEnterprise['partner_mobile'];
//        }

        // 只有门店导购才有门店权限
        $manageStores = [];
//        if (in_array($loginRole, [Constant::LOGIN_ROLE_STORE_MEMBER, Constant::LOGIN_ROLE_STORE_MANAGER])) {
//            try {
//                $dataRange = $account->dataRange([
//                    'admin_id' => $loginEnterprise['admin_id'],
//                    'client'   => Constant::APP_CLIENT_NAME
//                ]);
//                if (isset($dataRange['store_ids']) && !empty($dataRange['store_ids'])) {
//                    $manageStores = $dataRange['store_ids'];
//                }
//            } catch (Exception $exception) {
//                error($exception->getCode(), $exception->getMessage());
//            }
////        }
//        $newCache['mobile']                  = $loginEnterprise['mobile'];
//        $newCache['enterprise']              = array_values($memberInfoList);
//        $newCache['partner_enterprise_list'] = array_values($memberInfoList);
//        $newCache['choiceEnterprise']        = $memberInfoList[$requestData['enterprise_id']] ?? [];
////        $newCache['storeIds']                = $manageStores;
//        $newCache['enterprise_id']           = $requestData['enterprise_id'];
//        $newCache['role']                    = $loginRole;
//        $newCache['expiredTime']             = $cacheData['expiredTime'];
//        $newCache['loginId']                 = $cacheData['loginId'];

//        $expiredTime = Carbon::createFromTimestamp($cacheData['expiredTime']);
//        Cache::forget($requestData['login_cache_key']);
//        Cache::put($requestData['login_cache_key'], $newCache, $expiredTime);

        $enterprisePartnerStatus = self::PARTNER_STATUS_OFF;

        // 改为查micro-scrm的members表数据，storeworker 移除
//        $workers = (new StoreWorker())->getList([
//            'enterprise_id' => $requestData['enterprise_id'],
//            'mobile'        => $loginEnterprise['account_mobile'],
//            'is_page'       => false
//        ]);
//        $worker  = current($workers['list']);

        $memberInfo = GoMicroScrm::getMemberInfo([
            'enterprise_id' => $requestData['enterprise_id'],
            'mobile'        => $loginEnterprise['mobile'],
        ]);

        $authTargetUrl = '';
        if (empty($memberInfo['id'])) {
            error(652015);
        }

        $result = [
            'auth_flag'       => true,
            'auth_target_url' => $authTargetUrl,
            'store_id'        => $memberInfo['store_id'] ?? 0,
            'userid'          => $memberInfo['userid'] ?? '',
            'mem_uid'         => '',
            'mem_type'        => $loginUserType,
//            'partner_status'  => $enterprisePartnerStatus,
//            'version'         => $loginEnterprise['version'],
            'department_name' => $memberInfo['department_name'],
            'name'            => $memberInfo['name'],
        ];
        // 去除UC的查询
//        if (!empty($result['mem_uid'])) {
//            // 通过智慧导购的导购ID，查询导购的一级部门
//            $params   = [
//                'qy_domain' => $loginEnterprise['domain'],
//                'params'    => ['memUid' => $result['mem_uid']],
//            ];
//            $userInfo = (new Member())->detail($params);
//
//            $result['department_name'] = $userInfo['dpName'][0]['dpName'] ?? '';
//            $result['name']            = $userInfo['memUsername'] ?? '';
//        }

        return success($result);
    }


}
