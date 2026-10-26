<?php


namespace App\Http\Requests;


use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Services\EnterpriseAdminService;
use Illuminate\Validation\Rule;

class AccountRequest extends BaseValidateRequest
{
    /**
     * 创建账号验证
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function add($data)
    {
        $rule = [
            'operator_id'       => ['required'],
            'enterprise_id'     => ['required', 'integer'],
            'username'          => ['required'],
            'account'           => ['required', 'alpha_num'],
            'mobile'            => ['required'],
            //'password'          => ['required'],
            //'email'             => ['required'],
            'range_type'        => ['required', Rule::in(EnterpriseAdmin::RANGE_TYPE_ONLY, EnterpriseAdmin::RANGE_TYPE_ALL)],
            'range'             => ['required_if:range_type,' . EnterpriseAdmin::RANGE_TYPE_ONLY, 'array'],
            'range.*.type'      => ['required_if:range_type,' . EnterpriseAdmin::RANGE_TYPE_ONLY, 'integer'],
            'range.*.object_id' => ['required_if:range_type,' . EnterpriseAdmin::RANGE_TYPE_ONLY, 'integer'],
            'role_ids'          => ['required', 'array'],
            'role_ids.*'        => ['required', 'integer']
        ];

        $message = [
            'operator_id.required'          => '操作人ID不能为空',
            'enterprise_id.*'               => '企业ID不能为空',
            'username.required'             => '请输入姓名',
            'account.required'              => '请输入账号名称',
            'account.alpha_num'             => '账号名称必须是字母与数字的组合',
            'mobile.required'               => '请输入手机号',
            //'password.required'             => '请输入密码',
            //'email.required'                => '请输入邮箱',
            'range_type.required'           => '账号范围不存在或不是规定值',
            'range.required_if'             => '请选择账号范围',
            'range.array'                   => '账号范围数据类型错误',
            'range.*.type.required_if'      => '账号范围对象类型错误或不存在',
            'range.*.object_id.required_if' => '账号范围对象ID错误或不存在',
            'role_ids.required'             => '请选择账号角色',
            'role_ids.array'                => '账号所选角色数据类型错误',
            'role_ids.*.required'           => '请选择账号角色',
            'role_ids.*.integer'            => '所选角色数据类型错误',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 编辑账号验证
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function edit($data)
    {
        $rule = [
            'operator_id'            => ['required'],
            'enterprise_id'          => ['required', 'integer'],
            'username'               => ['required'],
            'account'                => ['required', 'alpha_num'],
            'password'               => ['required'],
            'mobile'                 => ['required'],
            //'email'                  => ['required'],
            'id'                     => ['required'],
            'range_type'             => ['required', Rule::in(EnterpriseAdmin::RANGE_TYPE_ONLY, EnterpriseAdmin::RANGE_TYPE_ALL)],
            'range'                  => ['required_if:range_type,' . EnterpriseAdmin::RANGE_TYPE_ONLY, 'array'],
            'range.*.type'           => ['required_if:range_type,' . EnterpriseAdmin::RANGE_TYPE_ONLY, 'integer'],
//            'range.*.object_id'      => ['required_if:range_type,' . EnterpriseAdmin::RANGE_TYPE_ONLY, 'integer'],
            'admin_available_system' => ['array']
        ];

        $message = [
            'operator_id.required'          => '操作人ID不能为空',
            'enterprise_id.*'               => '企业ID不能为空',
            'username.required'             => '请输入姓名',
            'account.required'              => '请输入账号名称',
            'account.alpha_num'             => '账号名称必须是字母与数字的组合',
            'password.required'             => '请输入密码',
            'mobile.required'               => '请输入手机号',
            //'email.required'                => '请输入email',
            'id.required'                   => '账号ID不能为空',
            'range_type.required'           => '账号范围不存在或不是规定值',
            'range.required_if'             => '请选择账号范围',
            'range.array'                   => '账号范围数据类型错误',
            'range.*.type.required_if'      => '账号范围对象类型错误或不存在',
            'range.*.object_id.required_if' => '账号范围对象ID错误或不存在',
            'admin_available_system.*'      => '可用系统传值有误',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 完善账号验证
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function perfect($data)
    {
        $rule = [
            'operator_id'   => ['required'],
            'enterprise_id' => ['required', 'integer'],
            'mobile'        => ['required'],
            'id'            => ['required'],
        ];

        $message = [
            'operator_id.required' => '操作人ID不能为空',
            'enterprise_id.*'      => '企业ID不能为空',
            'mobile.required'      => '手机号码不能为空',
            'id.required'          => '账号ID不能为空',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 账号详情
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function detail($data)
    {
        $rule = [
            'enterprise_id' => ['required', 'integer'],
            //'id'                => ['required'],
        ];

        $message = [
            'enterprise_id.*' => '企业ID不能为空',
            //'id.required'          => '管理员ID不能为空',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 删除账号验证
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function del($data)
    {
        $rule = [
            'enterprise_id' => ['required', 'integer'],
            'id'            => ['required'],
        ];

        $message = [
            'enterprise_id.*' => '企业ID不能为空',
            'id.required'     => '账号ID不能为空',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 修改账号状态
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function status($data)
    {
        $rule = [
            'enterprise_id' => ['required', 'integer'],
            'id'            => ['integer'],
            'status'        => ['required', Rule::in(EnterpriseAdmin::STATUS_OPEN, EnterpriseAdmin::STATUS_CLOSE)]
        ];

        $message = [
            'enterprise_id.*' => '企业ID不能为空',
            'id.integer'      => '账号ID不能为空',
            'status.*'        => '状态不在指定范围内或不存在',
        ];

        self::validator($data, $rule, $message);
    }

    /**
     * 账号列表验证
     *
     * @param $data
     *
     * @throws \Illuminate\Validation\ValidationException
     * @author lijian lijian@vchangyi.com 20202/04/13 10:55
     */
    public static function list($data)
    {
        $rule = [
            'enterprise_id' => ['required', 'integer'],
            'search_type'   => ['required'],
            'keyword'       => ['required_if:search_type,' . EnterpriseAdminService::SEARCH_MOBILE . ',' . EnterpriseAdminService::SEARCH_ACCOUNT]
        ];

        $message = [
            'enterprise_id.*'      => '企业ID不能为空',
            'search_type.required' => '搜索类型不能为空',
            'keyword.required_if'  => '搜索关键字不能为空',
        ];

        self::validator($data, $rule, $message);
    }

}