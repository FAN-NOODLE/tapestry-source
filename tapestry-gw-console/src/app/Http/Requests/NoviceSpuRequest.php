<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class NoviceSpuRequest
{
    /**
     * 新人优惠SPU添加验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/26 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'spu_id' => 'required',
            'limit_number' => 'required',
        ], [
            'spu_id.required' => '750101',
            'limit_number.required' => '750102'
        ]);
    }

    /**
     * 新人优惠SPU保存验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/26 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function save($data = [])
    {
        return Validator::make($data, [
            'novice_id' => 'required',
            'limit_number' => 'required',
        ], [
            'novice_id.required' => '750201',
            'limit_number.required' => '750102'
        ]);
    }

    /**
     * 新人优惠SPU状态修改验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/26 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function status($data = [])
    {
        return Validator::make($data, [
            'novice_id' => 'required',
            'status' => 'required',
        ], [
            'novice_id.required' => '750201',
            'status.required' => '750204'
        ]);
    }
}