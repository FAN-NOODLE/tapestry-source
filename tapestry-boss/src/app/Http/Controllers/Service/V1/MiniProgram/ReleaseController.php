<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:12
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

class ReleaseController extends MiniBaseController
{
    /**
     * 发布已通过审核的小程序
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index()
    {
        // 接收参数
        $versionId = request()->input('version_id');
        $enterpriseId = request()->input('enterprise_id');
        if (empty($versionId)) {
            return error(300110);
        }
        // 小程序app_id
        $appId = $this->appId;
        // 根据shop版本id获取指定审核id
        $auditInfo = $this->miniProgramService->getAuditIdByVersionId($versionId, $enterpriseId);
        // 发布指定审核id的版本
        $result = $this->miniProgramService->release($appId, $auditInfo['audit_id']);

        if (!$result) {
            return error(300108);
        }

        return success([
            'return_code'    => $result['errcode'],
            'return_message' => $result['errmsg']
        ]);
    }
}
