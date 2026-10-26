<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 15:34
 */

namespace App\Http\Controllers\Service\V1\Setting;


class AddController extends DataLakeConfigBase
{
    public function index()
    {
        $this->dataLakeConfigService->add($this->request->input());
        return success([]);
    }

    public function checkParam()
    {
        $rules    = [
            'enterprise_id' => 'required',
            'brand_code'    => 'required',
            'program_code'  => 'required',
//            'store_code'    => 'required',
//            'seller_name'   => 'required',
            'client_id'     => 'required',
            'client_secret' => 'required',
            'auth_user'     => 'required',
            'auth_password' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id缺失',
            'brand_code.required'    => '品牌编号必须存在',
            'program_code.required'  => '俱乐部编号必须存在',
            'client_id.required'     => 'httpBase 认证账号必须存在',
            'client_secret.required' => 'httpBase 认证密码必须存在',
            'auth_user.required'     => '获取token账号必须存在',
            'auth_password.required' => '获取token密码必须存在',

        ];
        $this->validate($this->request, $rules, $messages);
    }
}