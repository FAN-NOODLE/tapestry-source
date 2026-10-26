<?php

namespace App\Http\Requests\V1;

use App\Models\Enterprise;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;

class FastRegisterWeAppRequest extends BaseValidateRequest
{

    /**
     * 申请小程序
     *
     * @param array $data
     *
     * @throws ValidationException
     */
    public static function apply($data = [])
    {
        $rule = [
            'enterprise_id'        => 'required | integer',
            'legal_persona_name'   => 'required',
            'legal_persona_wechat' => 'required',
            'name'                 => 'required',
            'code_type'            => 'required | integer',
            'code'                 => 'required',
        ];

        $message = [
            'enterprise_id.required'        => '企业ID不能为空',
            'enterprise_id.integer'         => '企业ID类型错误',
            'legal_persona_name.required'   => '法人姓名不能为空',
            'legal_persona_wechat.required' => '法人微信不能为空',
            'name.required'                 => '企业名称不能为空',
            'code_type.required'            => '企业类型不能为空',
            'code_type.integer'             => '企业类型错误',
            'code.required'                 => '企业代码不能为空',
        ];
        self::validator($data, $rule, $message);
    }

}
