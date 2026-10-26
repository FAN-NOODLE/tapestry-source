<?php


namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class MaterialRequest
{
    /**
     * @param $params
     * @return \Illuminate\Contracts\Validation\Validator
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function add($params)
    {
        // 参数验证规则
        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'source_id'     => 'required|integer',
            'operate_id'    => 'required|integer',
            'spu_no'        => 'required|string',
            'name'          => 'required|string',
            'material_no'   => 'required|string',
            'cost_price'    => 'required|integer',
        ];

        // 提示语
        $messages = [
            'enterprise_id.*' => '企业ID不能为空',
            'source_id.*'     => '门店ID不能为空',
            'operate_id.*'    => '操作人不能为空',
            'name.*'          => '品名不能为空',
            'spu_no.*'        => '货品编号不能为空',
            'material_no.*'   => '物料编号不能为空',
            'cost_price.*'    => '成本价不能为空',
        ];

        return Validator::make($params, $rules, $messages);
    }

    /**
     * 编辑物料库存
     * @param $params
     * @return \Illuminate\Contracts\Validation\Validator
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function edit($params)
    {
        // 参数验证规则
        $rules = [
            'material_id' => 'required|integer|min:1',
            'cost_price'  => 'required|integer|min:0',
            'sign_price'  => 'required|integer|min:0',
            'name'        => 'required',
            'material_no' => 'required',
            'spu_no'      => 'required',
        ];

        // 提示语
        $messages = [
            'material_id.*' => '物料id不能为空',
            'cost_price.required'  => '成本价不能为空',
            'cost_price.integer'  => '成本价必须正整数',
            'cost_price.min'  => '成本价必须正整数',
            'name.*'        => '品名不能为空',
            'material_no.*' => '物料编码不能为空',
            'spu_no.*'      => '商品编码不能为空',
        ];

        validate($params, $rules, $messages);
    }

    /**
     * 更新物料库存验证
     * @param $params
     * @return \Illuminate\Contracts\Validation\Validator
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function updateStock($params)
    {
        // 参数验证规则
        $rules = [
            'material_id'   => 'required|integer|min:1',
            'num'           => 'required|integer|min:1',
            'operate_type'  => 'required|integer|in:1,2,3',
            'enterprise_id' => 'required|integer|min:1',
            'source_id'     => 'required|integer',
            'operate_id'    => 'required|integer',
        ];

        // 提示语
        $messages = [
            'material_id.*'   => '物料主键必填',
            'num.*'           => '修改数量必须正整数',
            'operate_type.*'  => '操作方式有误',
            'enterprise_id.*' => '企业ID不能为空',
            'source_id.*'     => '门店ID不能为空',
            'operate_id.*'    => '操作人不能为空',
        ];

        validate($params, $rules, $messages);

    }

    /**
     * 修改物料库存验证
     * @param $params
     * @return \Illuminate\Contracts\Validation\Validator
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function editStock($params)
    {
        // 参数验证规则
        $rules = [
            'material_id'   => 'required|integer|min:1',
            'material_no'   => 'required',
            'num'           => 'required|integer|min:1',
            'operate_type'  => 'required|integer|in:1,2,3',
        ];

        // 提示语
        $messages = [
            'material_id.*'   => '物料主键必填',
            'material_no.*'   => '物料编码必填',
            'num.*'           => '修改数量必须正整数',
            'operate_type.*'  => '操作方式有误',
        ];

        validate($params, $rules, $messages);
    }

}