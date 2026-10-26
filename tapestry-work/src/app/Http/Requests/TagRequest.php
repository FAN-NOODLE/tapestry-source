<?php


namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class TagRequest
{
    /**
     * @param $params
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function add($params)
    {
        $rules = [
            'qy_domain'        => 'required',
            'tag_list'         => 'required|array',
            'tag_list.*.name'  => 'required',
            'tag_list.*.order' => 'required',
        ];

        $message = [
            'qy_domain.required'        => '企业qy_domain不存在',
            'tag_list.required'         => '标签列表必传',
            'tag_list.array'            => '标签列表数据类型错误',
            'tag_list.*.name.required'  => '标签名称必传',
            'tag_list.*.order.required' => '标签排序必传,默认传0',
        ];

        $validator = Validator::make($params, $rules, $message);

        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }

        return true;
    }

    /**
     *
     * @param $params
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function edit($params)
    {
        $rules = [
            'qy_domain'    => 'required',
            'id'           => 'required',
            'name'         => 'required',
            'tag_group_id' => 'required',
            'type'         => 'required'
        ];

        $message = [
            'qy_domain.required'    => '企业qy_domain不存在',
            'id.required'           => '标签组或标签id必传',
            'name.required'         => '标签组或标签名称必传',
            'tag_group_id.required' => '标签分组id必传',
            'type.required'         => 'type类型必传'
        ];

        $validator = Validator::make($params, $rules, $message);

        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }

        return true;
    }

    /**
     * @param $params
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function tagRelation($params)
    {
        $rules = [
            'enterprise_id'    => 'required',
            'tag_id'           => 'required',
            'tag_group_id'     => 'required',
            'custom_group_ids' => 'array',
        ];

        $message = [
            'enterprise_id.required'    => '企业id不存在',
            'tag_id.required'           => '标签id必传',
            'tag_group_id.required'     => '标签组id必传',
            //'custom_group_ids.required' => '用户分组集合必传',
            'custom_group_ids.array'    => '用户分组集合数据类型错误',

        ];

        $validator = Validator::make($params, $rules, $message);

        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }

        return true;

    }

}