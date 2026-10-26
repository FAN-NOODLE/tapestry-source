<?php
/**
 * 下线小程序
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/11
 * Time: 16:11
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

class OfflineController extends MiniBaseController
{
    /***
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index()
    {
        $appId = $this->appId;
        if (!$appId) {
            return error(201101);
        }

        $res = $this->miniProgramService->changeOnlineStatus($appId, 'OFFLINE');

        return [
            'data' => [
                'status' => $res,
            ],
        ];
    }
}
