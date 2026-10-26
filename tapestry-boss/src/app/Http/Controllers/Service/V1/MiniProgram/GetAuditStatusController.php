<?php
/**
 *  查询小程序提审状态
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/11
 * Time: 15:29
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use Illuminate\Http\Request;

class GetAuditStatusController extends MiniBaseController
{
    /***
     * @param Request $request
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index(Request $request)
    {
        $appId = $this->appId;
        $auditId = $request->input('audit_id');

        if (!$appId) {
            return error(201101);
        }

        if (!$auditId) {
            return error(300002);
        }

        $auditLog = $this->miniProgramService->getAuditInfoByWXAuditId($appId, $auditId);

        return $auditLog;
    }
}
