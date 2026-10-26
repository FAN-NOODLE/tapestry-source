<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\App\V1;

use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use CrsSdk\Utils\Jaeger;
use Illuminate\Support\Facades\Crypt;

abstract class BaseController extends Controller
{
    /**
     * cookie key
     */
    const COOKIE_KEY = 'BOSS_ADMIN';

    /**
     * cookie 过期时间
     */
    const COOKIE_EXT = 72000;

    /**
     * 当前登录BOSS平台的管理员
     *
     * @var null|string
     */
    protected $adminId = null;

    /**
     * 是否必须登录
     *
     * @var bool
     */
    protected $requireLogin = true;

    /**
     * BaseController constructor.
     *
     * @throws \Exception
     */
    public function __construct()
    {
        parent::__construct();

        (new Jaeger())->serviceJaeger(
            'boss',
            [
                'path' => request()->path(),
            ]
        );

        // 非必须登录
        if (!$this->requireLogin) {
            return true;
        }
        $bossAdminId = $this->getCookie(self::COOKIE_KEY);

        if (empty($bossAdminId)) {
            return error(200004);
        }

        $this->adminId = $bossAdminId;

        return true;
    }

    /**
     * 设置cookie
     *
     * @param string $key
     * @param        $value
     *
     * @return bool
     * @author 陈朔  chenshuo@vchangyi.com
     * @date   2018/11/27 5:31 PM
     */
    protected function setCookie($key, $value)
    {
        setcookie($key, Crypt::encrypt($value), time() + self::COOKIE_EXT, '/');

        return true;
    }

    /**
     * 获取cookie
     *
     * @param string $key
     *
     * @return null|string
     * @author 陈朔  chenshuo@vchangyi.com
     * @date   2018/11/21 8:24 PM
     */
    private function getCookie($key)
    {
        if (!isset($_COOKIE[$key])) {
            throw new Error(100000, "请登录后操作");
        }
        return Crypt::decrypt($_COOKIE[$key]);
    }

    /**
     * @throws \Exception
     */
    public function __destruct()
    {
        (new Jaeger())->flush();
    }
}
