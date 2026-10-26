<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class EnterpriseAdminRequest
{
    /**
     * 企业管理员添加验证
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     * @param array $data 验证的数据
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required',
            'username'      => 'required',
            'account'       => 'required',
            'password'      => 'required'
        ], [
            'enterprise_id.required' => '企业ID不能为空',
            'username.required'      => '管理员名称不能为空',
            'account.required'       => '账户不能为空',
            'password.required'      => '密码不能为空'
        ]);
    }

    /**
     * 企业管理员密码修改
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/29 11:16
     *
     * @param array $data 验证的数据
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function passwordUpdate($data = [])
    {
        return Validator::make($data, [
            'admin_id'     => 'required',
            'old_password' => 'required',
            'password'     => 'required'
        ], [
            'admin_id.required'     => '企业管理员ID不能为空',
            'old_password.required' => '原密码不能为空',
            'password.required'     => '新密码不能为空'
        ]);
    }

    /**
     * 管理员删除验证
     *
     * @author zhaoqi
     * @date   2019-02-21
     *
     * @param array $data
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function delete($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required|integer',
            'admin_id'      => 'required|integer'
        ], [
            'enterprise_id.required' => 201101,
            'enterprise_id.integer'  => 201101,
            'admin_id.required'      => 202101,
            'admin_id.integer'       => 202101
        ]);
    }

    /**
     * 管理员状态更新验证
     *
     * @author zhaoqi
     * @date   2019-02-21
     *
     * @param array $data
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function status($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required|integer',
            'admin_id'      => 'required|integer',
            'status'        => 'required|integer'
        ], [
            'enterprise_id.required' => 201101,
            'enterprise_id.integer'  => 201101,
            'admin_id.required'      => 202101,
            'admin_id.integer'       => 202101,
            'status.required'        => 202111,
            'status.integer'         => 202111
        ]);
    }

    /**
     * 管理员编辑验证
     *
     * @author zhaoqi
     * @date   2019-02-21
     *
     * @param array $data
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function save($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required|integer',
            'admin_id'      => 'required|integer',
            'username'      => 'required',
            'account'       => 'required',
            'operator_id'   => 'required|integer',
            'role_id'       => 'required|integer|min:1'
        ], [
            'enterprise_id.required' => 201101,
            'enterprise_id.integer'  => 201101,
            'admin_id.required'      => 202101,
            'admin_id.integer'       => 202101,
            'username.required'      => 202107,
            'account.required'       => 202102,
            'operator_id.required'   => 202113,
            'operator_id.integer'    => 202113,
            'role_id.required'       => 202114,
            'role_id.integer'        => 202114,
            'role_id.min'            => 202114
        ]);
    }

    /**
     * 管理员编辑验证
     *
     * @author zhaoqi
     * @date   2019-02-21
     *
     * @param array $data
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function resetPassword($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required|integer',
            'admin_id'      => 'required|integer'
        ], [
            'enterprise_id.required' => 201101,
            'enterprise_id.integer'  => 201101,
            'admin_id.required'      => 202101,
            'admin_id.integer'       => 202101
        ]);
    }

    /**
     * 管理员添加 console 服务
     *
     * @author zhaoqi
     * @date   2019-02-21
     *
     * @param array $data
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function add($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required|integer',
            'username'      => 'required',
            'account'       => 'required',
            'password'      => 'required',
            'operator_id'   => 'required|integer',
            // todo::临时注释掉，门店权限规划好后去除注释
            // 'role_id'       => 'required|integer|min:1'
        ], [
            'enterprise_id.required' => 201101,
            'enterprise_id.integer'  => 201101,
            'username.required'      => 202107,
            'account.required'       => 202102,
            'password.required'      => 202112,
            'operator_id.required'   => 202113,
            'operator_id.integer'    => 202113,
            // todo::临时注释掉，门店权限规划好后去除注释
            // 'role_id.required'       => 202114,
            // 'role_id.integer'        => 202114,
            // 'role_id.min'            => 202114
        ]);
    }
}
