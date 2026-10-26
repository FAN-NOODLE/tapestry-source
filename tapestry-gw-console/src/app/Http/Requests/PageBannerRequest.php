<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class PageBannerRequest
{
    /**
     * 轮播图添加验证
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'id' => 'required',
            'picture' => 'required',
            'sort' => 'required',
        ], [
            'id.required' => '轮播图ID不能为空',
            'picture.required' => '图片不能为空',
            'sort.required' => '排序不能为空',
        ]);
    }
}