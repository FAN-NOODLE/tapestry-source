<?php


namespace App\Services;


use App\Jobs\MessageSendJob;
use App\Models\EnterpriseActivityMessage;
use App\Models\EnterpriseMessageScene;
use App\Models\SubscribeMessageRecord;
use App\Modules\Message\MessageFactory;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class ActivityMessageService extends BaseService
{
    /**
     * 发送活动开始提醒订阅
     * @return bool
     * @date 2021/1/12
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function SendMessage()
    {
        // 获取到达时间的记录
        $time = time();
        $messageConfigs = EnterpriseActivityMessage::query()
            ->where('send_time', '<=', $time)
            ->where('send_status', EnterpriseActivityMessage::SEND_STATUS_NO)
            ->where('enable_status', EnterpriseActivityMessage::ENABLE_STATUS_ON)
            ->orderBy('send_time')
            ->get()
            ->toArray();
        if (empty($messageConfigs)) {
            return true;
        }
        Log::info('activity-message:configs:', $messageConfigs);
        $enterpriseIds = array_column($messageConfigs, 'enterprise_id');
        $sceneIds = array_unique(array_column($messageConfigs, 'message_scene_id'));
        $sceneConfig = EnterpriseMessageScene::query()
            ->where('scene_type', EnterpriseActivityMessage::SCENE_TYPE)
            ->whereIn('message_scene_id', $sceneIds)
            ->whereIn('enterprise_id', $enterpriseIds)
            ->get()
            ->toArray();
        if (empty($sceneConfig)) {
            Log::debug('activity-message:not configs');
            return true;
        }
        $sceneConfig = array_column($sceneConfig, null, 'enterprise_id');

        // 获取用户记录
        $users = SubscribeMessageRecord::query()
            ->select(['user_id', 'object_type', 'id', 'enterprise_id'])
            ->whereIn('enterprise_id', $enterpriseIds)
            ->where('send_num', 0)
            ->where('fail_reason', '')
            ->where('scene_type', EnterpriseActivityMessage::SCENE_TYPE)
            ->orderBy('created_at')
            ->get()
            ->toArray();
        if (empty($users)) {
            $data = ['send_status' => EnterpriseActivityMessage::SEND_STATUS_YES];
            // 更新状态 为已发送
            EnterpriseActivityMessage::query()->whereIn('id', array_column($messageConfigs, 'id'))->update($data);
            return true;
        }
        $recordData = [];
        foreach ($users as $value) {
            $recordData[$value['enterprise_id']][] = $value;
        }

        foreach ($messageConfigs as $messageConfig) {
            if (!isset($sceneConfig[$messageConfig['enterprise_id']])) {
                Log::debug('activity-message:not config');
                continue;
            }
            // 停用 或 不可用 不进行发送
            if ($sceneConfig[$messageConfig['enterprise_id']]['status'] != 1) {
                // 更新状态 为停用
                EnterpriseActivityMessage::query()->where('id', $messageConfig['id'])->update(['enable_status' => EnterpriseActivityMessage::ENABLE_STATUS_OFF]);
                continue;
            }

            // 处理用户订阅记录，获取预计发送数量的人数
            $user = isset($recordData[$messageConfig['enterprise_id']]) ? array_splice($recordData[$messageConfig['enterprise_id']], 0, $messageConfig['expected_num']) : [];
            Log::debug('activity-message:users', $user);
            // 无订阅用户记录日志
            if (!empty($user)) {
                $userIds = array_column($user, 'user_id');
                $recordIds = array_column($user, 'id');
                // 关联记录
                SubscribeMessageRecord::query()
                    ->whereIn('id', $recordIds)
                    ->update(['object_id' => $messageConfig['id']]);
                // 发送消息
                $option = [
                    'page_url' => $messageConfig['page_url'],
                    'page_query' => $messageConfig['page_query'],
                    'type' => $messageConfig['message_scene_id'],
                    'user_ids' => $userIds,
                    'object_id' => $messageConfig['id'],
                    'object_type' => EnterpriseActivityMessage::OBJECT_TYPE,
                    'msg' =>  [
                        'live_title' => $messageConfig['title'],
                        'start_time' => $messageConfig['activity_start_time'],
                    ]
                ];
                $job = new MessageSendJob(
                    MessageFactory::MINI_PROGRAM,
                    $messageConfig['enterprise_id'],
                    $option['user_ids'],
                    EnterpriseActivityMessage::SCENE_TYPE,
                    $option['msg'],
                    $option);
                $job->handle();
            }

            $data = ['send_status' => EnterpriseActivityMessage::SEND_STATUS_YES];
            // 更新状态 为已发送
            EnterpriseActivityMessage::query()->where('id', $messageConfig['id'])->update($data);
        }


        return true;
    }

    /**
     * 同步订阅人数和发送人数
     * @return array
     * @date 2021/1/12
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function SyncMessageNumbers()
    {
        $type = EnterpriseMessageScene::query()
            ->select(['scene_type'])
            ->groupBy('scene_type')
            ->pluck('scene_type')
            ->toArray();
        $list = SubscribeMessageRecord::query()
            ->select(['enterprise_id', 'scene_type', DB::raw('count(id) as sub_num_total'), DB::raw('sum(send_num) as send_num_total')])
            ->whereIn('scene_type', $type)
            ->groupBy(['scene_type', 'enterprise_id'])
            ->get()
            ->toArray();
        Log::info("小程序统计人数数据", $list);
        foreach ($list as $item) {
            EnterpriseMessageScene::query()
                ->where('enterprise_id', $item['enterprise_id'])
                ->where('scene_type', $item['scene_type'])
                ->update([
                    'subscribe_num' => $item['sub_num_total'],
                    'send_num' => $item['sub_num_total'],
                    'send_real_num' => $item['send_num_total'],
                ]);
        }
        return $list;
    }
}