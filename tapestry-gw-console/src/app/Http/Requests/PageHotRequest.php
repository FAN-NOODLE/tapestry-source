<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class PageHotRequest
{
    /**
     * 爆品推荐添加验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'id' => 'integer',
            'picture' => 'required',
            'sort' => 'required',
        ], [
            'id.required' => '爆品推荐ID不能为空',
            'picture.required' => '图片不能为空',
            'sort.required' => '排序不能为空',
        ]);
    }
}