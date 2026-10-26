<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/22
 * Time: 3:05 PM
 */

namespace App\Http\Controllers\App\V1\Admin;

use App\Http\Controllers\App\V1\BaseController;

class LogoutController extends BaseController
{
    /**
     * 该接口无需强制登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 退出登录接口入口方法
     *
     * @return array|mixed
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/22 3:11 PM
     */
    public function index()
    {
        $result = setcookie(self::COOKIE_KEY, '', time()-7200, '/');

        if (isset($_COOKIE[self::COOKIE_KEY])) {
            unset($_COOKIE[self::COOKIE_KEY]);
        }

        if ($result) {
            return success(['data' => []]);
        }

        return error(100000);
    }
}
