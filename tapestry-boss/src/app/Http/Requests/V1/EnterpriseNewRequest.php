<?php

namespace App\Http\Requests\V1;

use App\Models\Enterprise;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;

class EnterpriseNewRequest extends BaseValidateRequest
{
    /**
     * 企业添加验证
     *
     * @param array $data
     * @throws ValidationException
     * @author 李坚 lijian@vchangyi.com
     * @date   2020/04/12 11:16
     *
     */
    public static function add($data = [])
    {
        $rule = [
            'name'                  => 'required',
            'contact'               => 'required',
            'mobile'                => 'required',
            'version'               => ['required', Rule::in([Enterprise::VERSION_CENTER_SHOP, Enterprise::VERSION_CENTER_SHOP_STORE, Enterprise::VERSION_DIRECT_CHAIN_STORE])],
            'account_limit'         => 'required | integer | min:0 | max:9999',
            'account_start_at'      => 'required | integer',
            // 'account_expire_at'     => 'required | integer',
            'service_merchant_id'   => 'required | integer',
            'email'                 => 'required | email',
            'password'              => 'required|min:15|confirmed|regex:/^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{15,}$/',
            'password_confirmation' => 'required|same:password',
            'industry_id'           => 'required',
            'corp_id'               => 'string',
        ];

        $message = [
            'name.required'        => '企业名称不能为空',
            'contact.required'     => '联系人不能为空',
            'mobile.required'      => '联系人电话不能为空',
            'industry_id.required' => '行业不能为空',

            'version.required' => '商城版本不能为空',
            'version.in'       => '商城版本编号不正确',

            'account_limit.required' => '最大子账号数不能为空',
            'account_limit.min'      => '最大子账号数不能小于0',
            'account_limit.max'      => '最大子账号数不能大于9999',
            'account_limit.integer'  => '最大子账号数类型错误',

            'account_start_at.required' => '账户有效期开始时间不能为空',
            'account_start_at.integer'  => '账户有效期开始时间类型错误',

            // 'account_expire_at.required' => '账户有效期到期时间不能为空',
            // 'account_expire_at.integer'  => '账户有效期到期时间类型错误',

            'email.required' => '超管账号（邮箱不能为空）',
            'email.email'    => '超管账号（邮箱）不是合法的邮箱格式',

            'password.required'  => '密码不能为空',
            'password.min'       => '密码最小15位',
            'password.confirmed' => '密码不一致',
            'password.regex'     => '密码格式错误',

            'service_merchant_id.required' => '支付服务商不能为空',
            'service_merchant_id.integer'  => '支付服务商类型错误',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 企业编辑验证
     *
     * @param array $data
     * @throws ValidationException
     * @author 李坚 lijian@vchangyi.com
     * @date   2020/04/12 11:16
     *
     */
    public static function edit($data = [])
    {
        $rule = [
            'name'              => 'required',
            'contact'           => 'required',
            'mobile'            => 'required',
            'version'           => ['required', Rule::in([Enterprise::VERSION_CENTER_SHOP, Enterprise::VERSION_CENTER_SHOP_STORE, Enterprise::VERSION_DIRECT_CHAIN_STORE])],
            'account_limit'     => 'required | integer | min:0 | max:9999',
            'account_start_at'  => 'required | integer',
            // 'account_expire_at' => 'required | integer',
            'enterprise_id'     => 'required',
            'available_system'  => 'array',
            'corp_id'           => 'string',
            'industry_id'       => 'required',
            'warehouse_type'    => 'required',
            'password'              => 'required|min:15|confirmed|regex:/^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{15,}$/',
            'password_confirmation' => 'required|same:password',
        ];

        $message = [
            'name.required'           => '企业名称不能为空',
            'contact.required'        => '联系人不能为空',
            'mobile.required'         => '联系人电话不能为空',
            'industry_id.required'    => '行业不能为空',
            'warehouse_type.required' => '企业仓库必选一个',
            'version.required'        => '商城版本不能为空',
            'version.in'              => '商城版本编号不正确',

            'account_limit.required' => '最大子账号数不能为空',
            'account_limit.min'      => '最大子账号数不能小于0',
            'account_limit.max'      => '最大子账号数不能大于9999',
            'account_limit.integer'  => '最大子账号数类型错误',

            'account_start_at.required' => '账户有效期开始时间不能为空',
            'account_start_at.integer'  => '账户有效期开始时间类型错误',

            // 'account_expire_at.required' => '账户有效期到期时间不能为空',
            // 'account_expire_at.integer'  => '账户有效期到期时间类型错误',

            'enterprise_id.required' => '企业ID不能为空',

            'password.required'  => '密码不能为空',
            'password.min'       => '密码最小15位',
            'password.confirmed' => '密码不一致',
            'password.regex'     => '密码格式错误',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 企业详情验证
     *
     * @param array $data
     * @throws ValidationException
     * @author 李坚 lijian@vchangyi.com
     * @date   2020/04/14 17:16
     *
     */
    public static function detail($data = [])
    {
        $rule = [
            'enterprise_id' => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业ID不能为空',
        ];

        self::validator($data, $rule, $message);
    }
}
