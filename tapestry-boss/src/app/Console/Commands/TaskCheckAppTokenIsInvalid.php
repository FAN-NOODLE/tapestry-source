<?php

namespace App\Console\Commands;

use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\EnterpriseNew\EnterpriseNewService;
use Illuminate\Console\Command;
use Overtrue\LaravelWeChat\Facade as EasyWeChat;

class TaskCheckAppTokenIsInvalid extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'task:check-app-token-is-invalid';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时任务检查小程序token是否失效';

    /**
     * TaskCheckAppTokenIsInvalid constructor.
     *
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        self::getAppList();
    }

    /**
     * 获取有效授权列表
     *
     * @return string
     */
    public function getAppList()
    {
        $result = EnterpriseThirdPartAccount::query()
            ->where([
                'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                'type'           => EnterpriseThirdPartAccount::TYPE_SHOP // 小程序类型
            ])->get()->toArray();
        $resultNew = array_column($result, null, 'app_id');
        $enterpriseIds = array_column($result, 'enterprise_id');


        $authorizerAppid = [];
        // 处理多开放平台
        $enterprise = Enterprise::query()
            ->distinct()
            ->whereIn('id', $enterpriseIds)
            ->get(['open_platform_id'])
            ->groupBy('open_platform_id');
        $openPlatformIds = array_unique(array_keys($enterprise->toArray()));

        $openPlatform = EasyWeChat::openPlatform();
        // 每个开放平台循环获取授权小程序信息，并输出到日志文件中
        foreach ($openPlatformIds as $openPlatformId) {
            changeOpenPlatformConfig(0, $openPlatform, $openPlatformId);
            do {
                $offset = $offset ?? 0;
                $appsList = $openPlatform->getAuthorizers($offset, 500);

                // 按格式输出整理后的数据
                foreach ($appsList['list'] as $appInfo) {
                    $authorizerAppid[] = $appInfo['authorizer_appid'];
                    if (isset($appInfo['refresh_token']) && !empty($appInfo['refresh_token']) && array_key_exists($appInfo['authorizer_appid'], $resultNew) && $resultNew[$appInfo['authorizer_appid']]['refresh_token'] != $appInfo['refresh_token']) {
                        file_put_contents('refresh_token.txt', 'app_id:'. $appInfo['authorizer_appid'] . '---token:'.$appInfo['refresh_token']. PHP_EOL, FILE_APPEND);
                        // 更新token
                        EnterpriseThirdPartAccount::query()
                            ->where('id', $resultNew[$appInfo['authorizer_appid']]['id'])
                            ->update(['refresh_token' => $appInfo['refresh_token']]);

                        $enterprise = Enterprise::query()
                            ->where('id', $resultNew[$appInfo['authorizer_appid']]['enterprise_id'])
                            ->first();

                        // 删除redis shop 授权数据缓存
                        delRedis([redisKey('SHOP_AUTH_DETAIL', ['enterprise_id' => $enterprise->id])]);
                        delRedis([redisKey('REDIS_SHOP_AUTH_PREFIX', ['hash' => $enterprise->hash])]);

                        (new EnterpriseNewService())->clearAuthInfoCache($enterprise->id);
                    }
                }

                $offset += 500;
                $total_count = $total_count ?? $appsList['total_count'];
            } while ($offset < $total_count);
        }

        $this->info('执行结果：-- 共计授权： '. $total_count);

        unset($offset, $total_count);

        return 'success';
    }
}
