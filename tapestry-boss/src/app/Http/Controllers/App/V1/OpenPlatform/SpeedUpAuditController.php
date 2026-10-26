<?php
/**
 * 开放平台：加急审核申请
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-10-31
 * Tome : 10:05
 */

namespace App\Http\Controllers\App\V1\OpenPlatform;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Library\OpenPlatform;
use App\Services\OpenPlatform\MiniProgramService;
use App\Services\OpenPlatform\OpenPlatformService;

class SpeedUpAuditController extends BaseController
{
    protected $requireLogin = false;
    /**
     * 审核id 即可，appId换auditId
     * @param MiniProgramService  $miniProgramService
     * @param OpenPlatformService $openPlatformService
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\HttpException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \EasyWeChat\Kernel\Exceptions\RuntimeException
     * @throws \GuzzleHttp\Exception\GuzzleException
     * @throws \Psr\SimpleCache\InvalidArgumentException
     */
    public function index(MiniProgramService $miniProgramService, OpenPlatformService $openPlatformService)
    {
        $appId  = request()->input('app_id');
        if (empty($appId)) {
            throw new Error(100000, '小程序app_id不存在');
        }
        // 获取代码审核id
        $auditInfo = $miniProgramService->recentAudit($appId);
        if (empty($auditInfo) || $auditInfo['status'] != 'AUDITING') {
            throw new Error(10000, '小程序审核信息不存在或状态不允许 ：' .$auditInfo['status'] ?? '');
        }
        $token  = $openPlatformService->getAuthorizeTokenByAppId($appId);
        $result = (new OpenPlatform())->speedUpAudit([
            'body'  => ['auditid' => $auditInfo['wx_audit_id']],
            'query' => ['access_token' => $token['authorizer_access_token']
            ]]);
        return success(['data' => $result]);
    }
}
