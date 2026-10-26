<?php

namespace App\Services;

use App\Common\RedisKey;
use App\Models\EnterpriseActivityMessage;
use App\Models\EnterpriseMessageScene;
use App\Models\SubscribeMessageRecord;
use App\Models\SubscribeMessageSetting;
use App\Models\SubscribeMessageTouch;
use App\Models\SubscribeMessageTouchTemplate;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class SubscribeMessageRecordService extends BaseService
{

    /**
     * 添加订阅记录
     *
     * @param $params
     * @param $priTmplId
     *
     * @return mixed
     */
    public function addSubscribeMessage($params, $priTmplId)
    {
        // 获取用户订阅状态
        $subscribeStatus = $this->getSubscribeStatus($params['object_id'], $params['user_id'], $params['object_type']);

        // 每个券都要进行通知
        if (SubscribeMessageSetting::SCENE_TYPE_COUPON_STATUS == $params['object_type']) {
            $subscribeStatus = SubscribeMessageRecord::SUBSCRIBE_FALSE;
        }

        // 处理上新订阅和活动开始提醒
        if (EnterpriseActivityMessage::OBJECT_TYPE == $params['object_type']) {
            $subscribeStatus = SubscribeMessageRecord::SUBSCRIBE_FALSE;
        }

        // 已订阅成功
        if ($subscribeStatus == SubscribeMessageRecord::SUBSCRIBE_TRUE) {
            $data['is_subscribe'] = SubscribeMessageRecord::SUBSCRIBED;
            return $data;
        }

        // 订阅消息模板ID
        $touchTemplate = SubscribeMessageTouch::query()
            ->where(['name' => $params['touch_id']])
            ->first();

        $objectTypeArray = [
            '25' => '15',
            '23' => '13',
            '24' => '14',
            '26' => '16',
            '28' => '16'
        ];

        $templates = SubscribeMessageTouchTemplate::query()
            ->whereIn('pri_tmpl_id', $priTmplId)
            ->when(in_array($params['object_type'], [16, 17]), function ($query) use ($priTmplId, $params, $objectTypeArray) {
                $recordTemplateIds = SubscribeMessageRecord::query()
                    ->whereIn('pri_tmpl_id', $priTmplId)
                    ->where('user_id', $params['user_id'])
                    ->whereIn('scene_type', array_keys($objectTypeArray))
                    ->where('send_num', 0)
                    ->where('fail_reason', '')
                    ->pluck('pri_tmpl_id')->all();
                return $query->whereNotIn('pri_tmpl_id', $recordTemplateIds);
            })
            ->get()->toArray();

        $templates = array_column($templates, null, 'pri_tmpl_id');

        $insert = [];
        $subscribeData = [];
        foreach ($priTmplId as $value) {
            if (in_array($params['object_type'], [16, 17]) && empty($templates[$value])) {
                continue;
            }
            $insert[] = [
                'user_id'       => $params['user_id'],
                'type'          => $params['type'],
                'object_type'   => $objectTypeArray[$templates[$value]['scene_type']] ?? $params['object_type'],
                'object_id'     => $params['object_id'],
                'touch_id'      => $touchTemplate['id'] ?? 0,
                'pri_tmpl_id'   => $value,
                'tid'           => $templates[$value]['tid'] ?? 0,
                'scene_type'    => $templates[$value]['scene_type'] ?? 0,
                'enterprise_id' => $templates[$value]['enterprise_id'] ?? 0,
                'created_at'    => time(),
            ];
            // 增加订阅数数据
            if ($templates[$value]['scene_type'] && $templates[$value]['enterprise_id']) {
                if (isset($subscribeData[$templates[$value]['scene_type']])) {
                    $subscribeData[$templates[$value]['scene_type']]['num'] += 1;
                } else {
                    $subscribeData[$templates[$value]['scene_type']]['num'] = 1;
                    $subscribeData[$templates[$value]['scene_type']]['enterprise_id'] = $templates[$value]['enterprise_id'];
                }
            }
        }

        try {
            SubscribeMessageRecord::query()->insert($insert);

            // 增加订阅数
            foreach ($subscribeData as $sceneType => $value) {
                EnterpriseMessageScene::query()
                    ->where('enterprise_id', $value['enterprise_id'])
                    ->where('scene_type', $sceneType)
                    ->update(['subscribe_num' => DB::raw('subscribe_num + 1' )]);
            }


            $this->redis->del([RedisKey::USER_SUBSCRIBE_NOTICE_STATUS . $params['user_id']]);

            $data['is_subscribe'] = SubscribeMessageRecord::SUBSCRIBE_TRUE;
        } catch (\Exception $e) {
            $data['is_subscribe'] = SubscribeMessageRecord::SUBSCRIBE_FALSE;
            Log::info('消息订阅失败' . '失败:参数' . json_encode($insert) . $e->getMessage());
        }

        return $data;
    }

    /**
     * 获取订阅状态
     *
     * @param $objectId
     * @param $userId
     * @param $objectType
     *
     * @return int
     */
    public function getSubscribeStatus($objectId, $userId, $objectType)
    {
        $subscribeInfo = $this->getSubscribeInfo([
            'user_id'     => $userId,
            'object_type' => $objectType,
            'object_id'   => $objectId
        ]);

        $subscribeStatus = SubscribeMessageRecord::SUBSCRIBE_TRUE;
        if (empty($subscribeInfo)) {
            $subscribeStatus = SubscribeMessageRecord::SUBSCRIBE_FALSE;
        }
        return $subscribeStatus;
    }

    /**
     * 获取订阅详情
     *
     * @param $params
     *
     * @return array
     */
    public function getSubscribeInfo($params)
    {
        $subscribeInfo = SubscribeMessageRecord::query()->where($params)
            ->where('send_num', '=', 0)
            ->where('fail_reason', '=', '')
            ->has('template')
            ->first();

        if (empty($subscribeInfo)) {
            return [];
        }
        return $subscribeInfo->toArray();
    }

    /**
     * 获取用户订阅列表
     *
     * @param $userId
     * @param $objectType
     *
     * @return array
     */
    public function getUserSubscribeList($userId, $objectType)
    {
        $subscribeList = SubscribeMessageRecord::query()
            ->where([
                'object_type' => $objectType,
                'user_id'     => $userId
            ])
            ->has('template')
            ->get()
            ->toArray();

        return $subscribeList ?? [];
    }

    /**
     * 获取用户订阅列表
     *
     * @param $params
     * @return array
     */
    public function getSubscribeList($params)
    {
        $size = $params['size'] ?? 1000;
        $size = $size <= 1000 ? $size : 1000;

        $subscribeList = SubscribeMessageRecord::query()
            ->select(['id', 'pri_tmpl_id', 'user_id', 'object_type', 'object_id', 'enterprise_id', 'scene_type', 'send_num', 'created_at'])
            ->where('send_num', '=', 0)
            ->where('fail_reason', '=', '')
            ->has('template')
            ->when($params['enterprise_id'] ?? 0, function ($query) use ($params) {
                return $query->where('enterprise_id', $params['enterprise_id']);
            })
            ->when($params['object_type'] ?? 0, function ($query) use ($params) {
                return $query->where('object_type', $params['object_type']);
            })
            ->when($params['object_id'] ?? 0, function ($query) use ($params) {
                return $query->where('object_id', $params['object_id']);
            })
            ->when($params['scene_type'] ?? 0, function ($query) use ($params) {
                return $query->where('scene_type', $params['scene_type']);
            })
            ->paginate($size);

        return ['list' => $subscribeList->items(),
                'meta' => [
                    'page'  => $subscribeList->currentPage(),
                    'size'  => $subscribeList->perPage(),
                    'total' => $subscribeList->total(),
                ]];
    }
}
