<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2020-05-06
 * Time: 09:31
 */

namespace App\Traits;

use App\Services\EnterpriseThirdPartAccountService;
use Illuminate\Support\Facades\Log;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;


trait OpenPlatformTrait
{
    /**
     * @param $enterpriseId
     *
     * @return \EasyWeChat\OpenPlatform\Authorizer\MiniProgram\Application
     */
    public function miniProgram($enterpriseId)
    {
        if (empty($enterpriseId)) {
            notice(700020);
        }

        $miniProgramInfo = (new EnterpriseThirdPartAccountService())->getMiniProgram($enterpriseId);

        if (empty($miniProgramInfo)) {
            notice(700021);
        }

        try {

            $openPlatform = LaravelWechat::openPlatform('default');
        } catch (\Exception $e) {
            Log::notice('实例化 open_platform 出错', ['code:' . $e->getCode(), 'message' . $e->getMessage()]);
            notice(800001);
        }
        // 处理兼容多开放平台
        changeOpenPlatformConfig($enterpriseId, $openPlatform);

        try {

            $appId        = $miniProgramInfo['app_id'];
            $refreshToken = $miniProgramInfo['refresh_token'];

            return $openPlatform->miniProgram($appId, $refreshToken);
        } catch (\Exception $e) {
            Log::notice('实例化 mini_program 出错', ['code:' . $e->getCode(), 'message' . $e->getMessage()]);
            notice(800002);
        }

    }
}