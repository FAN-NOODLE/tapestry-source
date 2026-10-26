<?php
/**
 * 下线小程序
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/11
 * Time: 16:11
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use Illuminate\Http\Request;

class OfflineController extends MiniBaseController
{
    /***
     * @param Request $request
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index(Request $request)
    {
        $appId = $request->input('app_id');
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
