<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/13
 * Time: 10:54 AM
 */

namespace App\Http\Controllers\Html\V1\Auth;

use App\Exceptions\Error;
use App\Http\Controllers\Html\V1\BaseController;
use App\Library\Uc\Common;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;

class LoginController extends BaseController
{
    /**
     * 微信授权地址
     *
     * @var string
     */
    const AUTH_URL = 'https://open.weixin.qq.com/connect/oauth2/authorize';

    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业微信html授权登录
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        if ($this->corpId == 'wx29da8c557a81167e') {
            file_put_contents(base_path('public/time.txt'), 'LoginController-1:' . msectime() . PHP_EOL, FILE_APPEND);
        }

        // 读取员工的cookie 信息
        $memberCookie = $this->getMemberCookie();

        $code = $request->post('code');

        // 存在员工cookie 已经确认身份
        if (!empty($memberCookie)) {

            if ($this->corpId == 'wx29da8c557a81167e') {
                file_put_contents(base_path('public/time.txt'), 'LoginController-2:' . msectime() . PHP_EOL, FILE_APPEND);
            }

            return self::success([
                'member' => $this->member,
                //'permission' => $permission['data']['permission'] ?? []
            ]);
        }

        // 没有code 没有 cookie 返回微信授权地址
        if (empty($code) && empty($memberCookie)) {

            $corpId      = $this->corpId;
            $agentId     = $this->agentId;
            $redirectUri = $request->post('redirect_uri');

            $authUrl = $this->getAuthUrl($corpId, $redirectUri, $agentId);

            if ($this->corpId == 'wx29da8c557a81167e') {
                file_put_contents(base_path('public/time.txt'), 'LoginController-3:' . msectime() . PHP_EOL, FILE_APPEND);
            }

            return self::success(['auth_url' => $authUrl]);
        }

        // 有code，根据code获取该员工（没有memUid 就是外部人员）的信息
        if (!empty($code)) {
            $memberData = (new Common())->authLogin($this->domain, $this->appIdentifier, ['code' => $code]);

            if ($this->corpId == 'wx29da8c557a81167e') {
                file_put_contents(base_path('public/time.txt'), 'LoginController-4:' . msectime() . '---requestId---' . $memberData['requestId'] ?? ''  . PHP_EOL, FILE_APPEND);
            }

            if (!empty($memberData['memUid']) && !empty($memberData['memUserid'])) {

                // 企业员工，内部人员
                $this->setMemberCookie(
                    Crypt::encrypt(json_encode(['memUid' => $memberData['memUid']]))
                );
                //$permission = $archService->member->memberPermission(['mem_uid' => $memberData['memUid']]);

                return self::success([
                    'member' => $memberData,
                    //'permission' => $permission['data']['permission'] ?? []
                ]);
            } else {
                throw new Error(1505, '对不起,仅限内部员工访问');
            }
        }

        return self::success([]);
    }

    /**
     * 获取微信授权地址
     *
     * @param string $corpId 企业微信唯一标识
     * @param string $redirectUri 回调地址
     * @param string $agentId 应用ID
     *
     * @return string
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/13 12:59 PM
     */
    protected function getAuthUrl($corpId, $redirectUri, $agentId)
    {
        $params = [
            'appid'         => $corpId,
            'redirect_uri'  => $redirectUri,
            'response_type' => 'code',
            'scope'         => 'snsapi_base',
            // 暂时弃用 by chenshuo
            // 'agentid' => $agentId,
            'state'         => uniqid()
        ];

        $queryString = http_build_query($params);

        $url = self::AUTH_URL . '?' . $queryString . '#wechat_redirect';
        return $url;
    }
}
