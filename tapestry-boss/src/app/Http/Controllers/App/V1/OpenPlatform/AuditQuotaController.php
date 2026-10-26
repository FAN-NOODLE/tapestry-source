<?php
/**
 * 开放平台：查询服务商的当月提审限额/加急 次数
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-10-31
 * Tome : 09:58
 */

namespace App\Http\Controllers\App\V1\OpenPlatform;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Library\OpenPlatform;
use App\Services\OpenPlatform\OpenPlatformService;

class AuditQuotaController extends BaseController
{
    protected $requireLogin = false;

    /**
     * 查询服务商的当月提审限额/加急 次数
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
    public function index(OpenPlatformService $openPlatformService)
    {
        $appId  = request()->input('app_id');
        if (empty($appId)) {
            throw new Error(100000, '小程序app_id不存在');
        }
        $token  = $openPlatformService->getAuthorizeTokenByAppId($appId);
        $result = (new OpenPlatform())->auditQuota([
            'query' => ['access_token' => $token['authorizer_access_token']
            ]]);
        return success(['data' => $result]);
    }
}
