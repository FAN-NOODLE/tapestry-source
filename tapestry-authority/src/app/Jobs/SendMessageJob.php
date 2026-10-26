<?php

namespace App\Jobs;

use App\Common\RedisKey;
use App\Models\EnterpriseActivityMessage;
use App\Models\EnterpriseMessageScene;
use App\Models\SubscribeMessageRecord;
use EasyWeChat\Factory;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

/**
 * 真正发送订阅消息的队列并开启多个消费者
 *
 * @package App\Jobs
 */
class SendMessageJob extends Job
{
    public $queue = 'send-message';

    private $config;
    private $recordID;
    private $content;
    private $redis;

    public function __construct($config, $recordID, $content)
    {
        $this->config = $config;
        $this->recordID = $recordID;
        $this->content = $content;
    }

    public function __wakeup()
    {
        $this->redis = redis();
    }

    public function handle()
    {
        Log::info('消息发送开始消费');

        $app = Factory::miniProgram($this->config);

        $record = SubscribeMessageRecord::query()->where('id', '=', $this->recordID)->first();
        if ($record === null) {
            return;
        }

        try {
            $res = $app->subscribe_message->send($this->content);
            $this->redis->del([RedisKey::USER_SUBSCRIBE_NOTICE_STATUS . $record->getAttribute('user_id')]);
        } catch (\Exception $exception) {
            Log::error('小程序订阅消息推送-失败', [
                'content' => $this->content,
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine()
            ]);
            $res = [
                'errcode' => 24503,
                'errmsg'  => $exception->getMessage(),
            ];
        }

        $recordUpdate = [];
        if (0 == $res['errcode']) {
            $result['success_list'][] = [
                'record'   => $record->getAttribute('id'),
                'response' => $res['errmsg']
            ];
            $recordUpdate['send_num'] = DB::raw('send_num + 1');
        } else {
            $result['fail_list'][] = [
                'record'   => $record->getAttribute('id'),
                'response' => $res['errmsg']
            ];
            $recordUpdate['fail_reason'] = $res['errmsg'];
        }
        $record->update($recordUpdate);

        // 更新发送数量
        $emsUpdate['send_num'] = DB::raw('send_num + 1');
        // 发送成功 更新真实发送数量
        if (isset($recordUpdate['send_num'])) {
            $emsUpdate['send_real_num'] = DB::raw('send_real_num + 1');
            if ($record->getAttribute('scene_type') == EnterpriseActivityMessage::SCENE_TYPE) {
                EnterpriseActivityMessage::query()
                    ->where('id', $record->getAttribute('object_id'))
                    ->update(['real_num' => DB::raw('real_num + 1')]);
            }
        }
        EnterpriseMessageScene::query()
            ->where('enterprise_id', $record->getAttribute('enterprise_id'))
            ->where('scene_type', $record->getAttribute('scene_type'))
            ->update($emsUpdate);


    }
}
