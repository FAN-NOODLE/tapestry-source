<?php

namespace App\Http\Middleware;

use App\Common\RedisKey;
use App\Models\EnterpriseMessageScene;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\MessageScene;
use Closure;

class CheckMessageStrategy
{

    protected $redis;

    public function __construct()
    {
        // ...
        $this->redis = redis();
    }

    /**
     * @param $request
     * @param Closure $next
     * @return mixed
     * @throws \Exception
     */
    public function handle($request, Closure $next)
    {
        $enterpriseId = $request->input('enterprise_id', 0);

        $enterprise = $this->getEnterpriseThirdPartAccountWithCategory($enterpriseId);

        if (empty($enterprise)) {
            return $next($request);
        }

        $configCount = count(config('message_strategy'));
        if ($configCount != $this->getEnterpriseMessageScenesCount($enterpriseId)) {
            $this->insertEnterpriseMessageScene($enterprise);
        }


        return $next($request);
    }

    public function insertEnterpriseMessageScene($enterprise)
    {
        $scenes = MessageScene::query()
            ->where('configure', 1)
            ->get()
            ->toArray();

        $enterpriseMessageScenes = EnterpriseMessageScene::query()
            ->where('enterprise_id', $enterprise['enterprise_id'])
            ->pluck('message_scene_id')
            ->all();
        $scenes = array_column($scenes, null, 'id');

        $messageStrategies     = config('message_strategy');
        $enterpriseCategoryIds = json_decode($enterprise['wx_category_ids'], 1);
        foreach ($messageStrategies as $key => &$item) {
            if (in_array($item['message_scene_id'], $enterpriseMessageScenes)) {
                unset($messageStrategies[$key]);
                continue;
            }

            if (!isset($scenes[$item['message_scene_id']])) {
                continue;
            }

            $availableCategoryIds = json_decode($scenes[$item['message_scene_id']]['available_category_ids'], 1);
            $intersectCategoryIds = array_intersect($availableCategoryIds, $enterpriseCategoryIds);

            // 1 启用  2 停用  3 行业不可用
            $item['status'] = 3;
            if ($intersectCategoryIds) {
                $item['status'] = 1;
            }
            $item['type'] = $item['type'] ?? 1;
            $item['enterprise_id'] = $enterprise['enterprise_id'];
            $item['created_at']    = time();
        }
        EnterpriseMessageScene::query()->insert($messageStrategies);
        $this->redis->del([RedisKey::ENTERPRISE_MESSAGE_SCENE_COUNT . $enterprise['enterprise_id']]);
        $this->redis->del([RedisKey::ENTERPRISE_MESSAGE_SCENE . $enterprise['enterprise_id']]);
    }


    protected function getEnterpriseThirdPartAccountWithCategory($enterpriseId)
    {
        $enterprise = $this->redis->get(RedisKey::ENTERPRISE_THIRD_PART_ACCOUNT_WX_CATEGORY . $enterpriseId);
        if ($enterprise) {
            $enterprise = json_decode($enterprise, 1);
        } else {
            $enterprise = EnterpriseThirdPartAccount::query()
                    ->select(['enterprise_id', 'app_id', 'app_secret', 'wx_category', 'wx_category_ids'])
                    ->where('enterprise_id', $enterpriseId)
                    ->where('app_id', '<>', '')
                    ->where('app_secret', '<>', '')
                    ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
                    ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
                    ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
                    ->where('wx_category_ids', '<>', '')
                    ->first() ?? [];

            if ($enterprise) {
                $enterprise = $enterprise->toArray();
            }
            $this->redis->set(RedisKey::ENTERPRISE_THIRD_PART_ACCOUNT_WX_CATEGORY . $enterpriseId, json_encode($enterprise, 320));
        }
        return $enterprise;
    }

    protected function getEnterpriseMessageScenesCount($enterpriseId)
    {
        $enterpriseMessageScenesCount = $this->redis->get(RedisKey::ENTERPRISE_MESSAGE_SCENE_COUNT . $enterpriseId);
        if (is_null($enterpriseMessageScenesCount)) {
            $enterpriseMessageScenesCount = EnterpriseMessageScene::query()
                ->where('enterprise_id', $enterpriseId)
                ->count();
            $this->redis->set(RedisKey::ENTERPRISE_MESSAGE_SCENE_COUNT . $enterpriseId, $enterpriseMessageScenesCount);
        }
        return $enterpriseMessageScenesCount;
    }
}
