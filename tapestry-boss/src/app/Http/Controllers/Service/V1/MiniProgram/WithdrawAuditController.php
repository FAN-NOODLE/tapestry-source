<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:11
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

class WithdrawAuditController extends MiniBaseController
{
    /**
     * 小程序撤审
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index()
    {
        // 小程序app_id
        $appId = $this->appId;
        if (!$appId) {
            return error(300100);
        }

        // 小程序审核撤回
        $result = $this->miniProgramService->cancelAudit($appId);
        if (!$result) {
            return error(300107);
        }

        return success([
            'return_code'    => $result['errcode'],
            'return_message' => $result['errmsg']
        ]);
    }
}
