<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class PageGuideRequest
{
    /**
     * 引导区添加验证
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
            'name' => 'required',
            'sort' => 'required',
            'view_column' => 'required',
        ], [
            'id.required' => '首页分类id不能为空',
            'picture.required' => '图片不能为空',
            'name.required' => '名称不能为空',
            'sort.required' => '排序不能为空',
            'view_column.required' => '商品展示列数不能为空',
        ]);
    }
}