<?php

namespace App\Console\Commands\Message;

use App\Common\RedisKey;
use App\Models\EnterpriseCategory;
use App\Models\EnterpriseThirdPartAccount;
use EasyWeChat\Factory;
use Illuminate\Console\Command;

class SyncThirdPartAccountCategory extends Command
{

    protected $signature = 'sync_third_part_account_category';

    protected $description = '小程序类目更新';

    protected $redis;

    public function __construct()
    {
        parent::__construct();
        $this->redis = redis();
    }

    /**
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function handle()
    {
        $this->updateThirdPartAccount();
        $this->insertCategory();
    }

    /**
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function updateThirdPartAccount()
    {
        $enterpriseList = EnterpriseThirdPartAccount::query()
            ->select(['id', 'enterprise_id', 'app_id', 'app_secret', 'name'])
            ->where('app_id', '<>', '')
            ->where('app_secret', '<>', '')
            ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->orderBy('enterprise_id', 'asc')
            ->get()->toArray();

        $enterpriseIds = array_values(array_column($enterpriseList, 'id'));

        $updates = [
            'wx_category'     => '',
            'wx_category_ids' => ''
        ];

        EnterpriseThirdPartAccount::query()
            ->whereNotIn('id', $enterpriseIds)
            ->update($updates);

        $config = config('wechat.defaults');

        foreach ($enterpriseList as $enterprise) {

            $this->info('企业' . $enterprise['enterprise_id'] . '-[' . $enterprise['name'] . ']正在处理');

            $config = array_merge($config, ['app_id' => $enterprise['app_id'], 'secret' => $enterprise['app_secret']]);

            $app = Factory::miniProgram($config);
            try {
                $category = $app->subscribe_message->getCategory();
            } catch (\Exception $exception) {
                $this->error($exception->getMessage());
            }

            if (!empty($category['data'])) {
                $updates = [
                    'wx_category'     => json_encode($category['data'], 320),
                    'wx_category_ids' => json_encode(array_column($category['data'], 'id'))
                ];
            }

            EnterpriseThirdPartAccount::query()
                ->where('enterprise_id', $enterprise['enterprise_id'])
                ->where('app_id', $enterprise['app_id'])
                ->where('app_secret', $enterprise['app_secret'])
                ->update($updates);
            $this->redis->del([RedisKey::ENTERPRISE_THIRD_PART_ACCOUNT_WX_CATEGORY . $enterprise['enterprise_id']]);
        }
    }

    /**
     *
     */
    public function insertCategory()
    {
        EnterpriseCategory::query()->where('name', '<>', '')->forceDelete();

        $enterpriseCategories = EnterpriseThirdPartAccount::query()
            ->where('wx_category_ids', '<>', '')
            ->pluck('wx_category')->all();

        foreach ($enterpriseCategories as &$enterpriseCategory) {
            $enterpriseCategory = json_decode($enterpriseCategory, 1);
        }
        $categories = collect($enterpriseCategories)->collapse()->unique()->sort()->values()->all();

        EnterpriseCategory::query()->insert($categories);
    }

    public function getCategoryIds()
    {
        $enterpriseCategoryIds = EnterpriseThirdPartAccount::query()
            ->where('wx_category_ids', '<>', '')
            ->pluck('wx_category_ids')->all();

        foreach ($enterpriseCategoryIds as &$enterpriseCategoryId) {
            $enterpriseCategoryId = json_decode($enterpriseCategoryId, 1);
        }

        return collect($enterpriseCategoryIds)->collapse()->unique()->sort()->values()->all();
    }
}
