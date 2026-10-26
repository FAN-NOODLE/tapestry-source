<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/22
 * Time: 3:18 PM
 */

namespace App\Http\Controllers\App\V1\Admin;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;

class MockLoginController extends BaseController
{
    /**
     * 该接口无需强制登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 模拟登录
     *
     * @param Request $request
     * @return array
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/22 3:21 PM
     */
    public function index(Request $request)
    {
        if (config('app.env') !== 'dev') {
            throw new Error(100000, '该接口禁止调用');
        }

        $mockAdminId = $request->get('id', 1);
        $this->setCookie(self::COOKIE_KEY, $mockAdminId);
        return success(['data' => []]);
    }
}
