<?php


namespace App\Http\Controllers\V1;


use App\Models\EnterpriseCategory;
use App\Models\ServiceMerchantSetting;
use App\Models\SubscribeMessageSetting;
use App\Models\SubscribeMessageTid;
use App\Models\SubscribeMessageTouch;

class TesterController extends BaseController
{
    public function setDefaultMerchant()
    {
        $id = $this->request->get('id');
        ServiceMerchantSetting::query()
            ->where('id', '!=', $id)
            ->update(['is_default' => 0]);
        ServiceMerchantSetting::query()
            ->where('id', $id)
            ->update(['is_default' => 1]);

        return success([]);
    }

    public function setMessageSetting()
    {
        $params = $this->request->only([
            'id',
            'category_id',
            'tid',
            'touch_ids',
            'scene_type',
            'scene_desc',
            'kids',
            'key_values',
            'page',
        ]);

        $this->validate($this->request, [
            'category_id' => 'required',
            'tid' => 'required',
            'touch_ids' => 'required | array',
            'scene_type' => 'required',
            'scene_desc' => 'required',
            'kids' => 'required | array',
            'key_values' => 'required | array',
        ], [
            'category_id.*' => '服务类目id必传',
            'tid.*' => '消息模板id必传',
            'touch_ids.*' => '触发事件id数组集合必传',
            'scene_type.*' => '场景类别必传',
            'scene_desc.*' => '服务场景描述必传',
            'kids.*' => '关键字下标集合必传',
            'key_values.*' => '模板参数必传',
        ]);
        // 验证服务类目
        $category = EnterpriseCategory::query()
            ->where('id', $params['category_id'])->exists();
        if (!$category) {
            return success(['服务类目不存在，请同步类目sync_third_part_account_category']);
        }
        // 验证消息模板
        $tid = SubscribeMessageTid::query()->where([
            ['categoryId', '=', $params['category_id']],
            ['tid', '=', $params['tid']],
            ['keywords', '!=', ''],
            ['rules', '!=', ''],
        ])->exists();
        if (!$tid) {
            return success(['模板关键字记录不存在，请同步关键字记录sync_message_tids']);
        }
        // 验证触发事件
        $touchCount = SubscribeMessageTouch::query()
            ->whereIn('id', $params['touch_ids'])->count();
        if ($touchCount != count($params['touch_ids'])) {
            return success(['触发事件有误，请确认触发事件是否存在']);
        }
        $data = [
            'category_id' => $params['category_id'],
            'tid' => $params['tid'],
            'touch_ids' => json_encode($params['touch_ids']),
            'scene_type' => $params['scene_type'],
            'scene_desc' => $params['scene_desc'],
            'kids' => json_encode($params['kids']),
            'key_values' => json_encode($params['key_values'], JSON_UNESCAPED_UNICODE + JSON_UNESCAPED_SLASHES, 512),
            'page' => $params['page'] ?? '',
        ];
        if (isset($params['id'])) {
            $setting = SubscribeMessageSetting::query()
                ->where('id', $params['id'])
                ->first();
            if (empty($setting)) {
                return success(['记录不存在']);
            }
            $setting->update($data);
        } else {
            $maxId = SubscribeMessageSetting::query()->max('id');
            $data['id'] = $maxId + 1;
            SubscribeMessageSetting::query()->create($data);
        }
        return success($data);
    }
}