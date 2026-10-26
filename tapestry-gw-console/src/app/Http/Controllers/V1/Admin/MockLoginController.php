<?php

namespace App\Http\Controllers\V1\Admin;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Crypt;
use App\Http\Controllers\V1\BaseController;

class MockLoginController extends BaseController
{
    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 模拟登录接口入口方法
     *
     * @return array
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/26 2:14 PM
     */
    public function index()
    {
        if (config('app.env') !== 'dev') {
            throw new Error(100000, '该接口仅仅用于开发环境');
        }

        $info = [
            'admin_id'      => 3,
            'enterprise_id' => 2,
            'serve_type'    => 0,
            'expire_time'   => time() + parent::TOKEN_EXPIRE * 21,
            'is_tmp' => 2
        ];
        $info['is_tmp'] = 2;
        $token = Crypt::encrypt(json_encode($info));

        return success([
            'data'  => $info,
            'token' => $token
        ]);
    }
}
