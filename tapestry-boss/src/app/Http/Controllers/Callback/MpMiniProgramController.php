<?php
/*
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Http\Controllers\Callback;

use App\Biz\OpenPlatformBiz;
use App\Events\MiniProgramAuthed;
use App\Exceptions\Error;
use App\Models\Enterprise;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class MpMiniProgramController extends BaseController
{
    const REDIS_SHOP_PREFIX = 'shopAuthInfo:';

    /**
     * 获取授权信息 (由微信调用).
     *
     * @param Request         $request
     * @param OpenPlatformBiz $openPlatformBiz
     * @param Enterprise      $enterprise
     *
     * @return \Illuminate\Http\RedirectResponse|\Illuminate\Routing\Redirector
     *
     * @throws \Exception
     *
     */
    public function index(Request $request, OpenPlatformBiz $openPlatformBiz, Enterprise $enterprise)
    {
        Log::info("authorize-access-token-params", $request->input());
        // 获取当前企业身份
        $enterprise = $enterprise->where('id', $request->input('enterprise_id'))->first();

        if (!$enterprise) {
            throw new  Error(100000, '企业不存在');
        }
        // 授权类型
        $authFor = $request->input('auth_for'); // SHOP;WORK
        $authCode = $request->input('auth_code');

        if (empty($authFor) || empty($authCode)) {
            throw new  Error(100000, '缺少必要参数');
        }

        // 记录授权情况w
        try {
            $thirdPartAccount = $openPlatformBiz->authorize($authCode, $enterprise->id, $authFor);
        } catch (\Exception $exception) {
            Log::error("authorize-access-token-error", [
                'message'   => $exception->getMessage(),
                'file'      => $exception->getFile(),
                'code'      => $exception->getCode(),
            ]);
            return redirect($request->input('back'));
        }

        if ('SHOP' == $authFor || 'WORK' == $authFor) {
            // 处理小程序授权完成事件
            event(new MiniProgramAuthed($thirdPartAccount));
        }

        // 删除redis shop 授权数据缓存
        delRedis([self::REDIS_SHOP_PREFIX . $enterprise->hash]);
        // 授权完成进行跳转
        $callbackUrl = $request->input('back');
        return redirect($callbackUrl);
    }
}
