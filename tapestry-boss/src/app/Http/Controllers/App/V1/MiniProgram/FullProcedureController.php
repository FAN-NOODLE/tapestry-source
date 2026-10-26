<?php
/**
 * 全网提审小程序
 * User: yingcai
 * Date: 2019/1/10
 * Time: 11:13 AM
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Http\Request;

class FullProcedureController extends BaseController
{
    protected $requireLogin = false;
    protected $miniProgramService;

    public function __construct(MiniProgramService $miniProgramService)
    {
        parent::__construct();
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
        $versionType = $request->input('version_type');
        $checked = (bool)$request->input('checked', false);
        if (!$versionType) {
            return error(300001);
        }
        $this->miniProgramService->fullProcedure($versionType, $checked, true);

        return success(['data' => []]);
    }
}
