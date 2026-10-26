<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-05-07
 * Time: 11:52
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

class RevertGrayReleaseController extends MiniBaseController
{
    /**
     * 小程序取消分阶段发布
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

        // 小程序取消分阶段发布
        $result = $this->miniProgramService->revertGrayRelease($appId);
        if (!$result) {
            return error(300110);
        }

        return success([
            'code'    => $result['errcode'],
            'message' => $result['errmsg']
        ]);
    }
}
