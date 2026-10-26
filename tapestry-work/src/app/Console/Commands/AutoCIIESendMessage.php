<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Models\ActivityReserveRecord;
use App\Models\SubscribeMessageRecord;
use App\Models\SubscribeMessageTemplate;
use App\Services\SubscribeMsgService;
use Exception;
use Illuminate\Console\Command;


class AutoCIIESendMessage extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:ciie:send_message';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'CIIE吊牌预约活动发送消息提醒';


    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     *
     * @throws Exception
     */
    public function handle()
    {
        $reserve_time = time() + 1800;
        // 查询预约活动
        $lists = ActivityReserveRecord::query()->where(['reserve_time' => $reserve_time, 'is_check' => 0])->get();

        // 为空
        if (empty($lists)) {
            return false;
        }

        foreach ($lists as $val) {
            $info = SubscribeMessageRecord::query()
                ->where('enterprise_id', $val['enterprise_id'])
                ->where('openid', $val['openid'])
                ->where('type', 1)
                ->where('send_num', 0)->where("rel_id", $val['id'])
                ->orderBy("id")
                ->first();
            if (empty($info)) {
                continue;
            }
            $tplInfo      = SubscribeMessageTemplate::query()
                ->where('enterprise_id', $val['enterprise_id'])
                ->where('pri_tmpl_id', $info->pri_tmpl_id)
                ->first();
            $dataNumArr   = explode(",", $tplInfo->data_num_str);
            $subscribeSvc = new SubscribeMsgService();
            list($appId, $appSecret) = $subscribeSvc->getMiniInfo($val['enterprise_id']);
            $delayTime = 3600;
            $msg       = [
                'thing' . $dataNumArr[0] => ['value' => "上海国家会展中心消费品展区6.2H展台"],//预约门店
                'thing' . $dataNumArr[1] => ['value' => "【Coach皮牌压印活动】"],//预约服务
                'time' . $dataNumArr[2]  => ['value' => date("Y-m-d", $val['reserve_time'])],//到店日期
                'time' . $dataNumArr[3]  => ['value' => date("H:i", $val['reserve_time']) . '~' . date("H:i", $val['reserve_time'] + $delayTime)],//到店时间
                'thing' . $dataNumArr[4] => ['value' => "请您携带相关电子凭证出席，期待您的莅临！"],//温馨提示
            ];
            $params    = [
                'touser'      => $val['openid'],
                'template_id' => $info->pri_tmpl_id,
                'page'        => $info->page,
                'data'        => $msg,
            ];

            $res = $subscribeSvc->send($appId, $appSecret, $params);
            if ($res['errcode'] !== 0) {
                $info->fail_reason = $res['errcode'] . '|' . $res['errmsg'];
            } else {
                $info->send_num++;
            }
            $info->save();
        }
        return 'ok';
    }
}