<?php

namespace App\Services;

use App\Common\RedisKey;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\SubscribeMessageSetting;
use App\Models\SubscribeMessageTouchTemplate;
use EasyWeChat\Factory;
use Illuminate\Support\Facades\Log;

class SubscribeMessageService extends BaseService
{
    public $enterprise;
    public $enterpriseId;
    public $app;


    public function __construct($enterpriseId)
    {
        parent::__construct();
        $this->enterpriseId = $enterpriseId;
        $this->enterprise = EnterpriseThirdPartAccount::query()
            ->select(['enterprise_id', 'app_id', 'app_secret', 'wx_category', 'wx_category_ids'])
            ->where('enterprise_id', $enterpriseId)
            ->where('app_id', '<>', '')
            ->where('app_secret', '<>', '')
            ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('wx_category_ids', '<>', '')
            ->first();

        $config = array_merge(config('wechat.defaults'), [
            'app_id' => $this->enterprise['app_id'],
            'secret' => $this->enterprise['app_secret']
        ]);

        $this->app = Factory::miniProgram($config);

        SubscribeMessageTouchTemplate::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->where('category_id', 0)
            ->forceDelete();
    }

    /**
     * @return array
     */
    public function sync()
    {
        Log::info('企业' . $this->enterpriseId . '正在处理');


        // 根据类目id获取配置列表
        $categoryIds = json_decode($this->enterprise['wx_category_ids'], 1);
        $settingList = SubscribeMessageSetting::query()
            ->whereIn('category_id', $categoryIds)
            ->where('key_values', '<>', '')
            ->get()->toArray();

        // 获取本地存储的个人模板列表
        $touchTemplateList = $this->getTouchTemplateList();

        // 同步配置项--需要增加的
        $needAddSettingIds = $this->addSettings($settingList, $touchTemplateList);

        // 同步配置项--需要删除的
        $needDelSettingIds = $this->delSettings($settingList, $touchTemplateList);

        // 同步配置项完成后重新查询数据库最新模板列表
        $touchTemplateListAfterSetting = $this->getTouchTemplateList();

        // 获取标记过tid_used的模板列表， tid_used 1=使用该模板  2=不使用该模板， tid_used=1才参与后续同步模板操作
        $touchTemplateListAfterTidUsed = $this->syncTidUsed($touchTemplateListAfterSetting);

        // 获取微信模板列表
        $wxTemplateList = $this->getWxTemplateList();

        // 微信模板 删除 (本地库不存在的模板)
        $delWxTmplIds = $this->delWxTmpls($wxTemplateList, $touchTemplateListAfterTidUsed);

        // 本地库 删除 （微信端不存在的模板）
        $delDbTmplIds = $this->delDbTmpls($wxTemplateList, $touchTemplateListAfterTidUsed);

        // 新增 （微信模板和本地都不存在的模板）
        $needAddTemplateIds = $this->addTemplates($settingList);

        if ($delWxTmplIds || $delDbTmplIds || $needAddTemplateIds) {
            Log::info("同步小程序模板消息--删除redis,企业id--" . $this->enterpriseId);
            $this->redis->del([RedisKey::SUBSCRIBE_TOUCH_DATA_CACHE . $this->enterpriseId]);
        }

        return [
            'needAddSettingIds'  => $needAddSettingIds,
            'needDelSettingIds'  => $needDelSettingIds,
            'delWxTmplIds'       => $delWxTmplIds,
            'delDbTmplIds'       => $delDbTmplIds,
            'needAddTemplateIds' => $needAddTemplateIds,
        ];
    }

    protected function getTouchTemplateList()
    {
        return SubscribeMessageTouchTemplate::query()
            ->select(['id', 'setting_id', 'tid', 'tid_used', 'touch_ids', 'pri_tmpl_id', 'scene_type', 'scene_desc', 'category_id'])
            ->where('enterprise_id', $this->enterpriseId)
            ->orderBy('scene_type', 'asc')
            ->orderBy('category_id', 'asc')
            ->get()->toArray();
    }

    protected function getWxTemplateList()
    {
        try {
            $wxTemplateList = $this->app->subscribe_message->getTemplates()['data'] ?? [];
        } catch (\Exception $exception) {
            $wxTemplateList = [];
            Log::error($exception->getMessage());
            Log::error('企业' . $this->enterpriseId . '小程序订阅消息同步模板id错误' . $exception->getMessage());
        }
        return $wxTemplateList;
    }

    protected function getTouchTemplateListAfterTidUses($touchTemplateList)
    {
        $collectSettings = collect($touchTemplateList);
        foreach ($touchTemplateList as &$needAddSetting) {
            $needAddSetting['tid_used'] = 2; // 1使用 2不用

            $sceneCount = $collectSettings
                ->where('scene_type', $needAddSetting['scene_type'])
                ->where('scene_desc', $needAddSetting['scene_desc'])
                ->count();

            if (1 == $sceneCount) {
                $needAddSetting['tid_used'] = 1;  // 1使用  2不用
            }

            if ($sceneCount > 1) {
                $minCategoryId = $collectSettings
                    ->where('scene_type', $needAddSetting['scene_type'])
                    ->where('scene_desc', $needAddSetting['scene_desc'])
                    ->min('category_id');
                if ($minCategoryId == $needAddSetting['category_id']) {
                    $needAddSetting['tid_used'] = 1;  // 1使用  2不用  用最小行业id的
                }
            }
        }
        return $touchTemplateList;
    }


    protected function syncTidUsed($touchTemplateList)
    {
        $touchTemplateListAfterTidUsed = $this->getTouchTemplateListAfterTidUses($touchTemplateList);
        $touchTemplateListAfterTidUsedCollect = collect($touchTemplateListAfterTidUsed);

        $syncUsedIds = $touchTemplateListAfterTidUsedCollect->where('tid_used', 1)->pluck('id')->values()->all();
        $syncNotUsedIds = $touchTemplateListAfterTidUsedCollect->where('tid_used', 2)->pluck('id')->values()->all();

        SubscribeMessageTouchTemplate::query()->whereIn('id', $syncUsedIds)->update(['tid_used' => 1]);
        SubscribeMessageTouchTemplate::query()->whereIn('id', $syncNotUsedIds)->update(['tid_used' => 2, 'pri_tmpl_id' => 1]);

        return $touchTemplateListAfterTidUsedCollect->where('tid_used', 1)->all();
    }


    protected function addSettings($settingList, $touchTemplateList)
    {
        $collectSettings = collect($settingList);
        $settingIds = array_values(array_column($settingList, 'id'));
        $templateSettingIds = array_values(array_column($touchTemplateList, 'setting_id'));

        $needAddSettingIds = array_values(array_diff($settingIds, $templateSettingIds));
        if ($needAddSettingIds) {
            Log::info('有需要添加的配置项' . json_encode($needAddSettingIds));
            $needAddSettings = $collectSettings->whereIn('id', $needAddSettingIds)->all();

            $createSettings = [];
            foreach ($needAddSettings as $needAddSetting) {
                $createSettings[] = [
                    'id'            => idGenerate(),
                    'setting_id'    => $needAddSetting['id'],
                    'enterprise_id' => $this->enterpriseId,
                    'tid'           => $needAddSetting['tid'],
                    'tid_used'      => $needAddSetting['tid_used'] ?? 2, // 1使用 2不用
                    'touch_ids'     => $needAddSetting['touch_ids'],
                    'category_id'   => $needAddSetting['category_id'],
                    'scene_type'    => $needAddSetting['scene_type'],
                    'scene_desc'    => $needAddSetting['scene_desc'],
                    'pri_tmpl_id'   => $needAddSetting['pri_tmpl_id'] ?? '',
                    'kids'          => $needAddSetting['kids'],
                    'created_at'    => time(),
                    'updated_at'    => time()
                ];
            };
            SubscribeMessageTouchTemplate::query()->insert($createSettings);
        }
        return $needAddSettingIds;
    }

    protected function delSettings($settingList, $touchTemplateList)
    {
        $settingIds = array_values(array_column($settingList, 'id'));
        $templateSettingIds = array_values(array_column($touchTemplateList, 'setting_id'));

        $needDelSettingIds = array_values(array_diff($templateSettingIds, $settingIds));
        if ($needDelSettingIds) {
            Log::info('有需要删除的配置项' . json_encode($needDelSettingIds));
            $needDelTemplates = collect($touchTemplateList)->whereIn('setting_id', $needDelSettingIds)->all();
            foreach ($needDelTemplates as $needDelTemplate) {
                $this->app->subscribe_message->deleteTemplate($needDelTemplate['pri_tmpl_id']);
            }
            SubscribeMessageTouchTemplate::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->whereIn('setting_id', $needDelSettingIds)
                ->delete();
        }
        return $needDelSettingIds;

    }

    protected function delWxTmpls($wxTemplateList, $touchTemplateList)
    {
        $wxTemplateIds = array_values(array_column($wxTemplateList, 'priTmplId'));
        $dbTemplateIds = array_values(array_column($touchTemplateList, 'pri_tmpl_id'));
        // 需删除的商户个人模板库模板id集合
        $delWxTmplIds = array_values(array_diff($wxTemplateIds, $dbTemplateIds));
        if ($delWxTmplIds) {
            Log::info('有需要删除的商户个人模板库模板' . json_encode($delWxTmplIds));

            foreach ($delWxTmplIds as $delWxTmplId) {
                $this->app->subscribe_message->deleteTemplate($delWxTmplId);
            }
        }
        return $delWxTmplIds;
    }

    protected function delDbTmpls($wxTemplateList, $touchTemplateList)
    {
        $wxTemplateIds = array_values(array_column($wxTemplateList, 'priTmplId'));
        $dbTemplateIds = array_values(array_column($touchTemplateList, 'pri_tmpl_id'));
        // 需删除的本地个人模板id集合
        $delDbTmplIds = array_values(array_diff($dbTemplateIds, $wxTemplateIds));
        if ($delDbTmplIds) {
            Log::info('有需要删除的数据库个人模板库模板' . json_encode($delDbTmplIds));

            SubscribeMessageTouchTemplate::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->whereIn('pri_tmpl_id', $delDbTmplIds)
                ->where('tid_used', 1)
                ->update(['pri_tmpl_id' => '']);
        }
        return $delDbTmplIds;
    }


    protected function addTemplates($settingList)
    {
        $needAddTemplates = SubscribeMessageTouchTemplate::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->where('tid_used', 1)
            ->where('pri_tmpl_id', '')
            ->get();

        $needAddTemplateIds = $needAddTemplates->pluck('setting_id')->all();

        if (count($needAddTemplates)) {
            Log::info('有需要添加的个人模板库模板settingIds' . json_encode($needAddTemplateIds));
            $settingList = array_column($settingList, null, 'id');
            foreach ($needAddTemplates as $needAddTemplate) {
                $kids = json_decode($settingList[$needAddTemplate['setting_id']]['kids'], 1);
                $sceneDesc = $settingList[$needAddTemplate['setting_id']]['scene_desc'];

                if (empty($kids)) {
                    continue;
                }
                try {
                    $res = $this->app->subscribe_message->addTemplate($needAddTemplate['tid'], $kids, $sceneDesc);
                } catch (\Exception $exception) {
                    Log::error($exception->getMessage());
                    $res = [
                        'errmsg' => 'service 503'
                    ];
                }

                $needAddTemplate->update([
                    'pri_tmpl_id' => $res['priTmplId'] ?? 1,
                    'fail_reason' => ('ok' == $res['errmsg']) ? '' : $res['errmsg']
                ]);
            }
        }
        return $needAddTemplateIds;
    }
}
