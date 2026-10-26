<?php
/*
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Http\Controllers\Callback;

use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Http\Request;

class MiniProgramController extends BaseController
{
    protected $miniProgramService;

    public function __construct(MiniProgramService $miniProgramService)
    {
        $this->miniProgramService = $miniProgramService;
    }

    /**
     * 自动上传小程序代码，并全网覆盖；返回错误未覆盖的小程序.
     *
     * @param Request $request
     *
     * @return string
     * @throws \App\Services\OpenPlatform\BaseException
     */
    public function index(Request $request)
    {
        // 防止请求执行超时
        set_time_limit(0);

        $versionType = $request->input('version_type');
        $checked = (bool)$request->input('checked', false);
        if (!$versionType) {
            return '小程序代码不存在';
        }
        $this->miniProgramService->fullProcedure($versionType, $checked);

        //event(new MiniProgram($versionType, $checked));
        return 'success';
    }
}
