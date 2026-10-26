<?php
/**
 * 直播间列表
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-11-05
 * Tome : 11:57
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use App\Library\OpenPlatform;
use App\Services\EnterpriseLogService;
use App\Services\OpenPlatform\OpenPlatformService;
use Illuminate\Support\Facades\Log;

class LiveInfoController extends MiniBaseController
{
    /**
     * @param OpenPlatformService $openPlatformService
     *
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\HttpException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \EasyWeChat\Kernel\Exceptions\RuntimeException
     * @throws \GuzzleHttp\Exception\GuzzleException
     * @throws \Psr\SimpleCache\InvalidArgumentException
     */
    public function index(OpenPlatformService $openPlatformService)
    {
        $page = $this->request->input('page', 1);
        $size = $this->request->input('size', 10);
        if (empty($this->appId)) {
            return error(201101);
        }
        try {
            $token  = $openPlatformService->getAuthorizeTokenByAppId($this->appId);
            $result = (new OpenPlatform())->getLiveInfo([
                'body'  => ['start' => $page, 'limit' => $size],
                'query' => ['access_token' => $token['authorizer_access_token']
                ]]);
            Log::info('app_id: '.$this->appId. ' 拉取直播房间列表 结果:' . json_encode($result));

            EnterpriseLogService::info('拉取直播房间列表结果', [
                'appId'  => $this->appId,
                'token'  => $token,
                'result' => $result
            ]);
        } catch (\Exception $e) {
            Log::info('app_id: '.$this->appId. ' 拉取直播房间列表 失败:' .$e->getMessage());
            EnterpriseLogService::error('拉取直播房间列表失败', [
                'appId'   => $this->appId,
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine()
            ]);
            $result = [];
        }
        $data = [];
        if (!empty($result) && $result['errcode'] == 0) {
            $data = $result['room_info'];
        }
        return success(['data' => $data]);
    }
}
