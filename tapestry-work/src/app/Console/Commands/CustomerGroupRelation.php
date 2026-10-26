<?php

namespace App\Console\Commands;

use App\Library\Cmq\CmqQueue;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class CustomerGroupRelation extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'customer_group_relation';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '每1小时从队列获取用户和分组绑定信息';

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
     * @throws \Exception
     */
    public function handle()
    {
       // 获取队列
        $queueName = config('cmq_setting.customer_group_relation');
        $queue = (new CmqQueue($queueName))->receiveMessage();

        if(!$queue) {
            Log::info('没有需要处理得用户分组关系队列数据');
            return true;
        }

        // 阅后即焚
        (new CmqQueue(config('cmq_setting.customer_group_relation')))->deleteMessage($queue->receiptHandle);

        $queueBody = json_decode($queue->msgBody, true);
        Log::info('用户分组关系：'.json_encode($queueBody));

        if ($queueBody['action'] != 'custom_group_relation') {
            Log::info('customer_group_relation：'.$queueBody['action']);
            return true;
        }

        // 处理数据 插入用户和用户分组关系表
        $query = $queueBody['query'];
        $relation = $query['relation'];
        $userIdS = array_column($relation, 'user_id');
        // 删除这些用户旧得分组数据
        \App\Models\CustomerGroupRelation::query()->withTrashed()->whereIn('user_id', $userIdS)->forceDelete();


        // 拼接数据批量插入
        $insertData = [];
        foreach ($relation as $value) {
            // 没有用户分组列表直接跳过
            if (empty($value['customer_groups'])) {
                continue;
            }
            foreach ($value['customer_groups'] as $v) {
                $insertData[] = [
                    'enterprise_id' => $query['enterprise_id'],
                    'customer_group_id' => $v,
                    'user_id' => $value['user_id'] ?? 0,
                    'ec_userid' => $value['ec_info'][0]['ec_userid'] ?? '',
                    'ec_info' => json_encode($value['ec_info']),
                    'receipt_handle' => $queue->receiptHandle
                ];
            }
        }

        Log::info('插入数据'.json_encode($insertData));
        \App\Models\CustomerGroupRelation::query()->insert($insertData);
        return true;
    }
}
