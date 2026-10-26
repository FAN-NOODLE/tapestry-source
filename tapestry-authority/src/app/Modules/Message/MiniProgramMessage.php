<?php

namespace App\Modules\Message;

use App\Common\RedisKey;
use App\Jobs\SendMessageJob;
use App\Models\SubscribeMessageSetting;
use App\Modules\Message\Content\MessageMapping;
use EasyWeChat\Factory;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Str;

class MiniProgramMessage extends AbstractMessage
{
    private $record;

    /**
     * @return array|mixed
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function send()
    {
        $result = [
            'success_list' => [],
            'fail_list'    => [],
            'fail_reason'  => ''
        ];

        if (null === $this->enterpriseThirdPartAccount) {
            $result['fail_reason'] = 'null enterprise';
            return $result;
        }

        if (null === $this->records) {
            $result['fail_reason'] = 'null records';
            return $result;
        }

        if ($this->enterpriseMessageScene && 1 <> $this->enterpriseMessageScene->status) {
            $result['fail_reason'] = 'status closed';
            return $result;
        }

        $config = array_merge(config('wechat.defaults'), [
            'app_id' => $this->enterpriseThirdPartAccount['app_id'],
            'secret' => $this->enterpriseThirdPartAccount['app_secret']
        ]);

        foreach ($this->records as $record) {
            // 用这个队列去消费，开启多个消费者
            $this->record = $record;
            dispatch(new SendMessageJob($config, $record['id'], $this->getContent()));
        }

        return $result;

        /*$app = Factory::miniProgram($config);

        foreach ($this->records as $record) {
            $this->record = $record;

            $content = $this->getContent();

            try {
                $res = $app->subscribe_message->send($content);
                $this->redis->del([RedisKey::USER_SUBSCRIBE_NOTICE_STATUS . $this->record['user_id']]);
                Log::info('小程序订阅消息推送', $content);
            } catch (\Exception $exception) {
                Log::error('小程序订阅消息推送-失败', [
                    'content' => $content,
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
                    'record'   => $record['id'],
                    'response' => $res['errmsg']
                ];
                $recordUpdate['send_num'] = DB::raw('send_num + 1');
            } else {
                $result['fail_list'][]       = [
                    'record'   => $record['id'],
                    'response' => $res['errmsg']
                ];
                $recordUpdate['fail_reason'] = $res['errmsg'];
            }
            $record->update($recordUpdate);
        }

        return $result;*/
    }

    /**
     * 获取组装好的消息内容
     *
     * @return array
     */
    protected function getContent(): array
    {
        return [
            'template_id' => $this->record['pri_tmpl_id'], // 所需下发的订阅模板id
            'touser'      => $this->getOpenid($this->record['user_id']),     // 接收者（用户）的 openid
            'page'        => $this->getPage(),
            'data'        => $this->getData(),
        ];
    }

    private function getPage()
    {
        $pageUrl   = empty($this->option['page_url']) ? $this->messageScene->page_url : $this->option['page_url'];
        $pageQuery = Str::After(($this->option['page_query'] ?? ''), '?');
        $cdpChan   = $this->messageScene->cdp_chan ?? '';

        if (1010 == $this->option['type']) {
            $cdpChan = (15 == $this->sceneType) ? 'kjsbtx' : 'kjcgtx';
        }

        $cdpQuery = 'cdp_chan1=dyxx&cdp_chan2=' . $cdpChan;

        if ($pageQuery && ($this->messageScene->cdp_params ?? '')) {
            parse_str($pageQuery, $res);
            $cdpQuery = $cdpQuery . ($res[$this->messageScene->cdp_params] ?? '');
        }

        return empty($pageQuery)
            ? $pageUrl . '?' . $cdpQuery
            : $pageUrl . $this->option['page_query'] . '&' . $cdpQuery;
    }

    private function getData()
    {
        $keyValues = json_decode($this->record['template']['setting']['key_values'], true);
        $sceneType = $this->record['template']['scene_type'];

        // 拼团进度提醒
        if (SubscribeMessageSetting::SCENE_TYPE_GROUP_PLAN == $sceneType) {
            $this->message['miss_num'] = $this->message['need_num'] - $this->message['group_progress'];
        }

        // 拼团结果通知
        if (SubscribeMessageSetting::SCENE_TYPE_GROUP_RESULT == $sceneType) {
            $this->message['notice'] = ('拼团成功' == $this->message['group_result']) ? '恭喜你拼团成功' : '很遗憾，拼团失败，再次发起试试';
        }

        // 好友助力
        if (SubscribeMessageSetting::SCENE_TYPE_HELP_COMPLETED == $sceneType) {
            $this->message['notice'] = (6 == $this->message['activity_type']) ? '恭喜您好友助力完成，请及时支付邮费' : '恭喜您好友助力完成，请及时查看';
        }

        $data = (new MessageMapping())->translate($keyValues, $this->message);

        return $data;
    }
}
