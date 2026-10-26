<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class CouponUserRequest
{
    /**
     * 礼劵领取用户列表验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/04 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function user($data = [])
    {
        return Validator::make($data, [
            'coupon_id' => 'required',
        ], [
            'coupon_id.required' => '礼劵ID不能为空',
        ]);
    }

    /**
     * 快到期礼券列表验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/04 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function expire($data = [])
    {
        return Validator::make($data, [
            'task_id' => 'required',
        ], [
            'task_id.required' => '任务ID不能为空',
        ]);
    }
}