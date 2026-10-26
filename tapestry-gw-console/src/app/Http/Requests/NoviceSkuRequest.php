<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class NoviceSkuRequest
{
    /**
     * 新人优惠SKU添加验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/26 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'sku_id' => 'required',
            'inventory_total' => 'required',
            'novice_price' => 'required',
            'status' => 'required',
            'is_enabled' => 'required',
        ], [
            'sku_id.required' => '750104',
            'inventory_total.required' => '750105',
            'novice_price.required' => '750106',
            'status.required' => '750107',
            'is_enabled.required' => '750111',
        ]);
    }

    /**
     * 新人优惠SKU编辑验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/26 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function save($data = [])
    {
        return Validator::make($data, [
            'sku_id' => 'required',
            'novice_sku_id' => 'required',
            'inventory_total' => 'required',
            'sales_total' => 'required',
            'novice_price' => 'required',
            'status' => 'required',
            'is_enabled' => 'required',
        ], [
            'sku_id.required' => '750104',
            'novice_sku_id.required' => '750202',
            'inventory_total.required' => '750105',
            'sales_total.required' => '750110',
            'novice_price.required' => '750106',
            'status.required' => '750107',
            'is_enabled.required' => '750111'
        ]);
    }
}