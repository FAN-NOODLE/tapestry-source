<?php
/**
 * 内购员工订单验证器规则
 */

namespace App\Http\Requests;

class StaffOrderRequest
{
    /**
     * 内购员工订单新增
     *
     * @param $params
     */
    public static function add($params)
    {
        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'order_no'      => 'required',
            'user_uid'       => 'required|integer',
            'order_money'   => 'required|integer',
        ];

        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',
            'order_no.required'      => '订单号不能为空',
            'user_uid.required'      => '用户ID不能为空',
            'user_uid.integer'       => '用户ID数据类型有误',
            'order_money.required'   => '订单金额数据不能为空',
            'order_money.integer'    => '订单金额数据类型有误',
        ];

        validator($params, $rules, $message);
    }

    /**
     * 内购员工订单更新
     *
     * @param $params
     */
    public static function update($params)
    {
        $rules = [
            'order_no'      => 'required',
        ];

        $message = [
            'order_no.required'      => '订单号不能为空',
        ];

        validator($params, $rules, $message);
    }

}
