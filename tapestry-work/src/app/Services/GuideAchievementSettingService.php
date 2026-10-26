<?php

namespace App\Services;

use App\Models\StaffSetting;

class GuideAchievementSettingService extends BaseService
{

    /**
     * 设置详情
     * @param $params
     * @return array
     */
    public function detail($params)
    {
        $detail = StaffSetting::query()
            ->when(!empty($params['enterprise_id']), function($query) use ($params) {
                $query->where('enterprise_id', $params['enterprise_id']);
            })
            ->where('type', StaffSetting::GUIDE_TYPE)
            ->first(['id', 'enterprise_id', 'binding_type', 'duration_time', 'change_owner', 'protected_time', 'commission_type', 'fission_type']);

        return $detail ? $detail->toArray() : [];
    }

    /**
     * 设置保存
     * @param $params
     * @return bool
     */
    public function save($params)
    {

        // 查询企业是否已有设置
        $setting = $this->detail(['enterprise_id' => $params['enterprise_id']]);
        if ($setting) {
            StaffSetting::query()
                ->where('id', $setting['id'])
                ->update($params);
        } else {
            StaffSetting::query()
                ->create($params);
        }

        // 清除缓存
        $guideSettingKey = redisKey('ENTERPRISE_GUIDE_SETTING', ['enterprise_id' => $params['enterprise_id'], 'type' => StaffSetting::GUIDE_TYPE]);
        delRedis([$guideSettingKey]);
        return true;
    }
}
