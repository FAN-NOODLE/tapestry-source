<?php


namespace App\Console\Commands\Message;


use App\Models\EnterpriseActivityMessage;
use App\Models\EnterpriseMessageScene;
use App\Models\MessageChannel;
use App\Models\MessageScene;
use App\Models\SubscribeMessageSetting;
use App\Models\SubscribeMessageTouch;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class InitMessageSceneAndChannel extends Command
{
    protected $signature = 'init_message_scene_and_channel {channel} {desc} {configure} {cdpChan} {--pageUrl=} {--touchScene=} {--messageScene=} {--defaultSendTime=} {--touchPage=} {--cdpParams=}';

    protected $description = '初始化消息场景和消息渠道 {mini_program|cdp} {消息场景} {0后台需要配置策略|1后台不需要配置策略} {小写首字母} {消息推送页面} {触发场景} {消息场景规则描述} {发送时间描述} {触发订阅页面} {cdp参数}';

    /**
     * 配置message_scene和message_channel表数据
     * @date 2021/1/19
     * @author 杨虎成 yanghucheng@vchangyi.com
     * 这两个基础表配置完成后 需要手动配置subscribe_message_setting
     */
    public function handle()
    {
        $channel   = $this->argument('channel');
        $desc      = $this->argument('desc');
        $configure = $this->argument('configure');
        $cdpChan   = $this->argument('cdpChan');

        $pageUrl         = $this->option('pageUrl');
        $touchScene      = $this->option('touchScene');
        $messageScene    = $this->option('messageScene');
        $defaultSendTime = $this->option('defaultSendTime');
        $cdpParams       = $this->option('cdpParams');
        $touchPage       = $this->option('touchPage');


        if (!in_array($channel, ['mini_program', 'cdp'])) {
            $this->info('channel错误！必须是mini_program或cdp');
            return;
        }
        if (!is_numeric($configure)) {
            $this->info('configure必须是数字0|1');
            return;
        } elseif (!in_array($configure, [0, 1])) {
            $this->info('configure必须是数字0|1');
            return;
        }
        if (empty($cdpChan)) {
            $this->info('cdpChan不能为空');
            return;
        }
        // 获取最大id
        $maxScene     = MessageScene::query()
            ->select([DB::raw('max(id) as id'), DB::raw('max(scene_type) as scene_type')])
            ->first()->toArray();
        $maxSceneSort = MessageScene::query()
            ->where('configure', 1)
            ->max('sort');

        $maxSceneId       = $maxScene['id'] + 1;
        $maxSceneType     = $maxScene['scene_type'] + 1;
        $messageSceneData = [
            'id'                => $maxSceneId,
            'scene_type'        => $maxSceneType,
            'desc'              => $desc,
            'page_url'          => $pageUrl ?? "",
            'touch_scene'       => $touchScene ?? "",
            'touch_page'        => $touchPage ?? "",
            'message_scene'     => $messageScene ?? "",
            'default_send_time' => $defaultSendTime ?? "",
            'cdp_chan'          => $cdpChan,
            'configure'         => $configure,
            'cdp_params'        => $cdpParams ?? "",
            'sort'              => $configure == 1 ? $maxSceneSort + 1 : 0
        ];

        $maxChannelId = MessageChannel::query()->max('id');
        $channelData  = [
            'id'               => $maxChannelId + 1,
            'message_scene_id' => $maxSceneId,
            'scene_type'       => $maxSceneType,
            'channel'          => $channel,
            'status'           => 1
        ];

        try {
            DB::beginTransaction();
            // 1. message_scene
            MessageScene::query()->create($messageSceneData);

            // 2. message_channel
            MessageChannel::query()->create($channelData);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            $this->info($e->getMessage());
            return;
        }
        $this->info('success');
    }
}