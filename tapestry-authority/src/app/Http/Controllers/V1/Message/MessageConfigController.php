<?php


namespace App\Http\Controllers\V1\Message;


use App\Http\Controllers\V1\BaseController;
use App\Models\MessageChannel;
use App\Models\MessageScene;
use App\Models\SubscribeMessageSetting;
use App\Models\SubscribeMessageTouch;
use Illuminate\Support\Facades\Artisan;
use Illuminate\Support\Facades\DB;

class MessageConfigController extends BaseController
{
    public function index()
    {
        $channel              = $this->request->input('channel', 'mini_program');
        $isCopy               = $this->request->input('is_copy', 0);
        $copyMessageSceneType = $this->request->input('copy_message_scene_type', 0);
        if ($isCopy && empty($copyMessageSceneType)) {
            error(100000, [
                'message' => 'copy_message_scene_type不能为空'
            ]);
        }
        $this->validate($this->request, [
            'desc'              => 'required', // 描述
            'page_url'          => 'required', // 地址
            'touch_scene'       => 'required', // 触发场景
            'touch_page'        => 'required|array', // 触发页面，例：["首页","购物车"]
            'message_scene'     => 'required', // 消息发送场景
            'default_send_time' => 'required', // 默认发送时间
            'cdp_chan'          => 'required',
            'touch_ids'         => 'array', // 触发事件表id数组
            'touch_name'        => 'string', // 触发事件名称（前端提供）
            'comment'           => 'string' // 触发事件描述
        ]);
        $params = $this->request->all();
        // 获取最大id
        $maxId = MessageScene::query()
            ->where('configure', 1)
            ->max('id');

        $maxType = MessageScene::query()
            ->max('scene_type');

        $sceneId          = $maxId + 1;
        $sceneType        = $maxType + 1;
        $messageSceneData = [
            'id'                => $sceneId,
            'scene_type'        => $sceneType,
            'desc'              => $params['desc'],
            'page_url'          => $params['page_url'],
            'touch_scene'       => $params['touch_scene'],
            'touch_page'        => json_encode($params['touch_page'], JSON_UNESCAPED_UNICODE),
            'message_scene'     => $params['message_scene'],
            'default_send_time' => $params['default_send_time'],
            'cdp_chan'          => $params['cdp_chan'],
            'configure'         => $params['configure']
        ];
        $maxChannelId     = MessageChannel::query()->max('id');
        $channelData      = [
            'id'               => $maxChannelId + 1,
            'message_scene_id' => $sceneId,
            'scene_type'       => $sceneType,
            'channel'          => $channel,
            'status'           => 1
        ];

        $touchData = false;
        // 验证是否传递 触发事件表id数组
        if (isset($params['touch_ids']) && !empty($params['touch_ids'])) {
            $touchIds = SubscribeMessageTouch::query()
                ->whereIn('id', $params['touch_ids'])
                ->get()
                ->toArray();
            if (empty($touchIds)) {
                error(100000, ['message' => '触发事件不存在']);
            }
            if (count($touchIds) != count($params['touch_ids'])) {
                error(100000, ['message' => '触发事件不不匹配']);
            }
            $touchIds = array_column($touchIds, 'id');
        } else {
            if (empty($params['touch_name'])) {
                error(100000, ['message' => 'touch_name不能为空']);
            }
            if (empty($params['comment'])) {
                error(100000, ['message' => 'comment不能为空']);
            }
            $maxTouchId = SubscribeMessageTouch::query()->max('id');
            $touchId    = $maxTouchId + 1;
            $touchData  = [
                'name'    => $params['touch_name'],
                'comment' => $params['comment'],
                'id'      => $touchId
            ];
            $touchIds   = [$touchId];
        }

        $settingData  = [];
        $maxSettingId = SubscribeMessageSetting::query()
            ->max('id');
        if ($isCopy) {
            $copyList = SubscribeMessageSetting::query()
                ->where('scene_type', $copyMessageSceneType)
                ->get()
                ->toArray();
            foreach ($copyList as $item) {
                $maxSettingId  += 1;
                $settingData[] = [
                    'id'          => $maxSettingId,
                    'category_id' => $item['category_id'],
                    'tid'         => $item['tid'],
                    'type'        => $item['type'],
                    'kids'        => $item['kids'],
                    'key_rules'   => $item['key_rules'],
                    'key_values'  => $item['key_values'],
                    'scene_type'  => $sceneType,
                    'scene_desc'  => $params['desc'],
                    'touch_ids'   => json_encode($touchIds)
                ];
            }
        }

        try {
            DB::beginTransaction();
            // 1. message_scene
            MessageScene::query()->create($messageSceneData);

            // 2. message_channel
            MessageChannel::query()->create($channelData);

            // 3. touch
            if ($touchData) {
                SubscribeMessageTouch::query()->create($touchData);
            }
            // 4. setting
            if (!empty($settingData)) {
                SubscribeMessageSetting::query()->insert($settingData);
            }


            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            error(100000, [
                'message' => $e->getMessage()
            ]);
        }
        $scenes = MessageScene::query()
            ->with('channels.settings.category:id,name')
            ->orderBy('id', 'asc')
            ->get();

        foreach ($scenes as $scene) {
            $channel = $scene['channels']->firstWhere('channel', 'mini_program');
            if (is_null($channel)) {
                continue;
            }

            $categoryIds = $channel['settings']->pluck('category_id')->sort()->values()->all();
            $categories  = $channel['settings']->pluck('category')->sortBy('id')->values()->all();

            MessageScene::query()
                ->where('id', $scene['id'])
                ->update([
                    'available_category_ids' => json_encode($categoryIds),
                    'available_categories'   => json_encode($categories, 320)
                ]);
        }
        return success([
            'scene'   => $messageSceneData,
            'channel' => $channelData,
            'setting' => $settingData
        ]);
    }
}