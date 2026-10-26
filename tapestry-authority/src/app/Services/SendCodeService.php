<?php

namespace App\Services;

use App\Models\EnterpriseAdmin;

class SendCodeService extends BaseService
{
    /**
     * 验证修改手机号时点击获取验证码判断要修改的手机号是否占用
     * @param $rule
     * @param $mobile
     * @return bool
     */
    public function sendCodeRule($rule, $mobile)
    {
        // 此规则只有在管理员完善信息修改自己手机号时生效
        if (isset($rule['check_enterprise_admin'])) {
            // 查询传入的手机号是否能查询到其他企业-能查到不允许修改
            $noEnterprise = EnterpriseAdmin::query()
                ->where('account', 'like', '%'.$mobile.'%')
                ->get()
                ->toArray();
            if (!empty($noEnterprise)) {
                notice(701007); // 不可修改为此手机号 已经在其他企业使用
            }

            return true;
        }
    }

}