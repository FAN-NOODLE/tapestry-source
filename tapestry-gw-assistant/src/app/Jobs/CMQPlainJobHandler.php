<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-15
 * Time: 10:51
 */

namespace App\Jobs;

use App\Events\StoreOrderStatus;
use Illuminate\Contracts\Queue\Job as QueueJob;
use Illuminate\Support\Facades\Log;

class CMQPlainJobHandler
{
    /**
     * 队列名和处理事件映射关系
     *
     * @var array
     */
    private $eventMap
        = [
            'store-order-status' => StoreOrderStatus::class
        ];

    /**
     * @param QueueJob $job
     * @param $payload
     */
    public function handle(QueueJob $job, $payload)
    {
        try {
            // 获取队列名称
            $queueName   = $job->getQueue();
            $queueData   = json_decode($payload, true);
            $environment = env('APP_ENV');
            $startTime   = microtime(true);
            Log::info("[assistant] {$startTime} 接收到{$queueName}消息", $queueData);
            if (preg_match("/{$environment}-([a-zA-z\-_0-9]+)/", $queueName, $match)) {
                $realQueueName = $match[1];
                // 触发事件
                event(new $this->eventMap[$realQueueName]($queueData, $startTime));
            }

            // 删除消息
            if (!$job->isDeletedOrReleased()) {
                $job->delete();
            }
        } catch (\Exception $exception) {
            Log::error("消息处理失败，{$exception->getMessage()}");
        }

    }
}
