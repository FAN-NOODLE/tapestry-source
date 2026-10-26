<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class CouponRequest
{
    /**
     * 礼劵添加验证
     *
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'name'            => 'required',
            'discount_amount' => 'required',
            'condition_value' => 'required',
            'start_time'      => 'required',
            'end_time'        => 'required',
        ], [
            'name.required'            => '礼劵名称不能为空',
            'discount_amount.required' => '优惠金额不能为空',
            'condition_value.required' => '使用条件不能为空',
            'start_time.required'      => '开始时间不能为空',
            'end_time.required'        => '结束时间不能为空',
        ]);
    }
}
