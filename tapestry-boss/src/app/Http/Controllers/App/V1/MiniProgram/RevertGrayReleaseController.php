<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-05-07
 * Time: 11:52
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Http\Request;

class RevertGrayReleaseController extends BaseController
{
    /**
     * 小程序取消分阶段发布
     *
     * @param Request            $request
     * @param MiniProgramService $miniProgramService
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

        // 小程序取消分阶段发布
        $result = $miniProgramService->revertGrayRelease($appId);
        if (!$result) {
            return error(300110);
        }

        return success([
            'code'    => $result['errcode'],
            'message' => $result['errmsg']
        ]);
    }
}
