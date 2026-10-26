<?php

namespace App\Services;

use App\Models\MobileWhite;
use Illuminate\Http\JsonResponse;

/**
 * 导购身份在weclient功能身份的白名单
 */
class MobileWhiteService
{
    /**
     * @param $enterpriseId
     * @param $mobile
     * @param int $type
     * @return bool
     */
    public static function WhiteList($enterpriseId, $mobile, int $type = MobileWhite::CUSTOMER_SEARCH):bool
    {
        // 查询是否不验证
        $where['type'] = $type;
        $where['enterprise_id'] = $enterpriseId;
        $exists = MobileWhite::query()
            ->where($where)
            ->where("status", MobileWhite::ALL_ALLOW) //存在2则全部不验证
            ->exists();
        if ($exists) {
            //如果不验证则直接返回
            return true;
        }
        //查询当前用户是否有权限
        $mobileWhiteLists = MobileWhite::query()
            ->where($where)->pluck('mobile')->toArray();
        return in_array($mobile, $mobileWhiteLists);
    }
}