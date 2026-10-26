<?php


namespace App\Http\Requests;

use http\Exception\InvalidArgumentException;
use Illuminate\Support\Facades\Validator;

class CustomerGroupRequest
{
    /**
     * 员工内购运营设置表
     *
     * @param $params
     *
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function add($params)
    {
        $rules = [
            'enterprise_id'    => 'required|integer|min:1',
            'name'  => 'required',
            'user_visit'   => 'required',
            'coupon'      => 'required',
            'order' => 'required',
        ];

        $message = [
            'enterprise_id.required'    => '企业ID不存在',
            'enterprise_id.integer'     => '企业ID数据类型有误',
            'enterprise_id.min'         => '企业ID数据类型有误',

            'name.required'  => '用户分组名称不能为空',

            'user_visit.required'   => '浏览记录必传',

            'coupon.required'      => '账户优惠券必传',

            'order.required' => '购买行为必传',

        ];

        $validator = Validator::make($params, $rules, $message);

        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }

        return true;
    }

}