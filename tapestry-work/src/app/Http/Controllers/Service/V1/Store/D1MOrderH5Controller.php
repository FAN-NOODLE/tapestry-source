<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Models\StoreDepartment;

class D1MOrderH5Controller extends BaseController
{

    const TYPE_BA_ORDER = 1;  //导购订单

    const TYPE_SHARE_ORDER = 2;  //导购分享

    /**
     * 获取D1M订单跳转url
     * @return array
     * @throws \Exception
     */
    public function index(): array
    {
        $params = $this->request->only([
            'enterprise_id',
            'mobile',
            'type',
            'openid'
        ]);
        if (empty($params['enterprise_id']) || empty($params['mobile']) || empty($params['type'])) {
            throw new Error(600000, '参数错误');
        }
        $memberInfo = Members::query()->with('departmentsRel.departments')
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('mobile', $params['mobile'])
            ->first(['id', 'userid', 'position', 'type']);

        if (empty($memberInfo) || empty($memberInfo->departmentsRel)) {
            return success();
        }
        $memberInfo = $memberInfo->toArray();


        $dpInfo = current($memberInfo['departments_rel']);
        $storeInfo = StoreDepartment::query()->where('enterprise_id', $params['enterprise_id'])
            ->where('wechat_department_id', $dpInfo['departments']['third_id'])
            ->first();
        if (empty($storeInfo)) {
            return success(['url' => '']); //如果未开启在线订单 则直接返回空
        }
        $storeInfo = $storeInfo->toArray();
        if ($storeInfo['popup_store_status'] != StoreDepartment::POPUP_STORE_OPEN) {
            return success(['url' => '']); //如果未开启在线订单 则直接返回空
        }
        if ($params['type'] == self::TYPE_BA_ORDER) {
            $url = $storeInfo['sa_check_order_url'] . '?';
        } else {
            $url = $storeInfo['sa_share_order_url'] . '?';
        }
        $rolename = 'sa';
        if ($memberInfo['type'] == Members::STORE_SM
            || in_array($memberInfo['position'], Members::COACH_STORE_SM)
            || ($storeInfo['sa_check_role_status'] == StoreDepartment::SA_CHECK_ROLE_BY_OTHERS
                && in_array($memberInfo['userid'], explode(',', $storeInfo['sa_check_role'])))
        ) {
            //如果是企业管理、职位是店长、开启了其他人也可以管理且在管理范围则设置为sm
            $rolename = 'sm';
        }
        $url .= 'empcode=' . $memberInfo['userid'] . '&';
        $url .= 'shopcode=' . $storeInfo['third_store_code'] . '&';
        $url .= 'rolename=' . $rolename;
        if ($params['type'] == self::TYPE_SHARE_ORDER && !empty($params['openid'])) {
            $userList = WxUserCorpExternal::getExternalUserInfoByOpenId($params['enterprise_id'], [$params['openid']]);
            $userInfo = current($userList);
//            $url .= '&memberid=' . $userInfo['unionid'] . '&';
            $url .= '&openid=' . $userInfo['unionid'] . '&';
            $url .= 'nickname=' . urlencode($userInfo['nick']);
        }
        //导购订单
        return success(['url' => $url]);

    }
}
