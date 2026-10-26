<?php
namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class AdminRequest
{
    /**
     * 企业管理员密码修改
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/29 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function passwordUpdate($data = [])
    {
        return Validator::make($data,[
            'old_password' => 'required',
            'password' => 'required',
            're_password' => 'required'
        ],[
            'old_password.required' => '原密码不能为空',
            'password.required' => '新密码不能为空',
            're_password.required' => '确认密码不能为空'
        ]);
    }
}