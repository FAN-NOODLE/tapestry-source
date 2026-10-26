<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-09-27
 * Tome : 15:06
 */

namespace App\Console\Commands;

use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;
use App\Models\EnterpriseThirdPartAccount;
use App\Library\SubscribeMessageSdk;

class SyncMiniCategory extends Command
{


    /**
     * @var string
     */
    protected $signature = 'auto:sync_mini_category';

    /**
     * @var string
     */
    protected $description = '更新小程序分类';

    public function handle()
    {
        Log::info('更新小程序分类, 开始');
        $list = app(MiniProgramService::class)->getMiniBaseList(['void_wx_category' => true]);

        foreach ($list as $item) {
            try {
                // 获取当前小程序分类
                $getCategoryData = (new SubscribeMessageSdk($item['enterprise_id'], $item['app_id'], $item['app_secret']))->getCategory();

                if ($getCategoryData['errcode'] == 0) {
                    $wxCategoryIds = json_encode(array_column($getCategoryData['data'], 'id'));
                    EnterpriseThirdPartAccount::query()->where('app_id', $item['app_id'])->update(
                        [
                            'wx_category_ids' => $wxCategoryIds,
                            'wx_category'     => json_encode($getCategoryData['data'])
                        ]
                    );
                }
            }catch (\Exception $exception){
                Log::info('报错企业'.json_encode($item));
            }
        }

        Log::info('更新小程序分类, 完成');
    }
}
