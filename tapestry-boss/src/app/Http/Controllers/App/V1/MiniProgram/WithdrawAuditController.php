<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:11
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Http\Request;

class WithdrawAuditController extends BaseController
{
    /**
     * 小程序撤审
     *
     * @param Request            $request
     * @param MiniProgramService $miniProgramService 小程序相关服务
     *
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index(Request $request, MiniProgramService $miniProgramService)
    {
        // 小程序app_id
        $appId = $request->post('app_id', '');
        if (!$appId) {
            return error(300100);
        }

        // 小程序审核撤回
        $result = $miniProgramService->cancelAudit($appId);
        if (!$result) {
            return error(300107);
        }

        return success([
            'code'    => $result['errcode'],
            'message' => $result['errmsg']
        ]);
    }
}
