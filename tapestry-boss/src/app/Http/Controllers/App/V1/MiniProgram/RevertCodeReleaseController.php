<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-05-06
 * Time: 14:57
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Http\Request;

class RevertCodeReleaseController extends BaseController
{
    /**
     * 小程序版本回退
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

        // 小程序版本回退
        $result = $miniProgramService->rollbackRelease($appId);
        if (!$result) {
            return error(300109);
        }

        return success([
            'code'    => $result['errcode'],
            'message' => $result['errmsg']
        ]);
    }
}
