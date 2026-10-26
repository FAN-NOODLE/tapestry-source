<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/19
 * Time: 11:20 AM
 */

namespace App\Http\Controllers\Html\V1\Auth;

use App\Exceptions\Error;
use App\Http\Controllers\Html\V1\BaseController;
use App\Library\Uc\JsApi;
use Illuminate\Http\Request;

class AgentConfigController extends BaseController
{
    /**
     * 不需要登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * js config
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws BaseException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/19 11:21 AM
     */
    public function index(Request $request)
    {
        if ($this->corpId == 'wx29da8c557a81167e') {
            file_put_contents(base_path('public/time.txt'), 'AgentConfigController-1:' . msectime() . PHP_EOL, FILE_APPEND);
        }
        $url = $request->input('url');

        if (empty($url)) {
            throw new Error(2000, '参数错误');
        }

        $params = [
            'url' => $url,
            'time' => time(),
            'noncestr' => uniqid()
        ];

        try {
            $data = (new JsApi())->agentConfig($this->domain, $this->appIdentifier, $params);
            if ($this->corpId == 'wx29da8c557a81167e') {
                file_put_contents(base_path('public/time.txt'), 'AgentConfigController-2:' . msectime() . '---requestId---' . $data['requestId'] ?? ''  . PHP_EOL, FILE_APPEND);
            }
            return self::success($data);
        } catch (\Exception $exception) {
            throw new Error(1000, 'JS API config 签名生成错误'.$exception->getMessage());
        }
    }
}
