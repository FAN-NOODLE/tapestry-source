<?php
/**
 * HelpRequest.php
 * Created on 2019-09-30 10:50
 * Created by shizhixiong@vchangyi.com
 */

namespace App\Http\Requests;

use App\Models\StaffRight;
use App\Models\StaffSetting;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;

class StaffRightRequest
{
    /**
     * 下载中心导出新增
     *
     * @param $params
     *
     */
    public static function add($params)
    {
        $rules = [
            'enterprise_id'                  => 'required|integer|min:1',
            'qy_domain'                      => 'required',
            'staff_buy_list.*.obj_id'        => 'required',
            'staff_buy_list.*.obj_parent_id' => 'required',
            'staff_buy_list.*.obj_type'      => ['required', Rule::in(StaffRight::TYPE_STAFF, StaffRight::TYPE_DEP)],
            'setting.binding_type'           => ['required', Rule::in(StaffSetting::FOREVER, StaffSetting::EXPIRE)],
            'setting.change_owner'           => ['required', Rule::in(StaffSetting::NOT_ALLOW, StaffSetting::ALLOW)],
            'setting.commission_type'        => ['required', Rule::in(StaffSetting::HAS_COMMISSION, StaffSetting::NOT_HAS_COMMISSION)],
            'setting.settle_type'            => ['required', Rule::in(StaffSetting::TRADE_SETTLE, StaffSetting::POST_SALE_SETTLE)],
        ];

        $message = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',

            'qy_domain.required' => '企业domain不存在',

            'staff_buy_list.*.obj_id.required' => '员工或部门id不能为空',

            'staff_buy_list.*.obj_parent_id.required' => '员工或部门父级id不能为空',

            'staff_buy_list.*.obj_type.required' => '员工或部门类型不存在',
            'staff_buy_list.*.obj_type.in'       => '员工类型为1，部门类型为2',

            'setting.required' => '员工设置不存在',
            'setting.array'    => '员工设置数据类型有误',

            'staff_buy_list.binding_type.required' => '绑定时效不能为空',
            'staff_buy_list.binding_type.in'       => '绑定时效数据有误',

            'staff_buy_list.change_owner.required' => '是否允许抢客不能为空',
            'staff_buy_list.change_owner.in'       => '是否允许抢客数据有误',

            'staff_buy_list.commission_type.required' => '内购员工购买商品是否有佣金不能为空',
            'staff_buy_list.commission_type.in'       => '内购员工购买商品是否有佣金数据有误',

            'staff_buy_list.settle_type.required' => '佣金结算方式不能为空',
            'staff_buy_list.settle_type.in'       => '佣金结算方式数据有误',
        ];

        validator($params, $rules, $message);
    }
}
