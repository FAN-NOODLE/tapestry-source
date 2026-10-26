<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/4
 * Time: 5:34 PM
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;

class MockController extends BaseController
{
    /**
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 模拟登录
     *
     * @param Request $request
     * @return array
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/4 5:40 PM
     */
    public function index(Request $request)
    {
        if (config('app.env') === 'prod') {
            throw new Error(10000, '小伙子，你走错路了');
        }

        $params = $request->only(['memUid', 'corpId', 'sessionKey', 'qyDomain', 'userId']);

        // 如果有就模拟，没有就采用一个默认的
        if (!empty($params) && is_array($params)) {

            if (empty($params['corpId']) || empty($params['memUid']) || empty($params['qyDomain'] || empty($params['userId']))) {
                throw new Error(100000, '参数错误，请检查 memUid corpId qyDomain userId 是否传值正确');
            }

            $token = $params;
            // 前端没法模拟sessionKey 这个要小心，和sessionKey相关必须真机调试，不能模拟！！！
            $token['sessionKey'] = 'mock';
        } else {
            // 拼接token信息
            $token = [
                'memUid' => 'B76C515B7F000001280AE5402C67E265',
                'corpId' => 'wwd4f9da8b88237b47',
                'sessionKey' => '',
                'qyDomain' => '2E80C2677F000001341C116FA85B6EC7',
                'userId' => 'huhansan'
            ];
        }

        $tokenData = Crypt::encrypt(json_encode($token));

        return success(
            [
                'data' => [
                    'token' => $tokenData,
                    'mock' => $token
                ]
            ]
        );
    }
}
