<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/13
 * Time: 1:18 PM
 */

namespace App\Http\Controllers\Html\V1\Auth;

use App\Exceptions\Error;
use App\Http\Controllers\Html\V1\BaseController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;

class MockController extends BaseController
{
    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业微信手机端模拟登录
     *
     * @param Request $request
     * @return mixed
     * @throws BaseException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/13 1:21 PM
     */
    public function index(Request $request)
    {
        /*
         * 生产环境禁止模拟登录
         * 开发/测试环境允许模拟登陆
         */
        if (config('app.env') == 'prod' ) {
            throw new Error(401,'无访问权限，该接口仅仅用于开发阶段');
        }

        $memUid = $request->get('mem_uid');

        if (empty($memUid)) {
            throw new Error(2000, '模拟登录的人员不能为空');
        }

        // 先清除一次cookie
        setcookie(self::COOKIE_MEMBER_NAME, 0, time() - 72000, '/');

        // 模拟内部人员
        if (!empty($memUid)) {
            $this->setMemberCookie(Crypt::encrypt(json_encode(['memUid' => $memUid])));
            return self::success([$this->member]);
        }

        throw new Error(2000, '模拟登录失败了！');
    }
}