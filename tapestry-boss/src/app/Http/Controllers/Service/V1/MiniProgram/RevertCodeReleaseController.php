<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-05-06
 * Time: 14:57
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

class RevertCodeReleaseController extends MiniBaseController
{
    /**
     * 小程序版本回退
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

        // 小程序版本回退
        $result = $this->miniProgramService->rollbackRelease($appId);
        if (!$result) {
            return error(300109);
        }

        return success([
            'return_code'    => $result['errcode'],
            'return_message' => $result['errmsg']
        ]);
    }
}
