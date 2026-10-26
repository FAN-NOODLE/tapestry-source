<?php

namespace App\Sdk\ShenCe;

use BatchConsumer;
use DebugConsumer;
use Exception;
use FileConsumer;
use SensorsAnalytics;

class Client
{
    const FileConsumer  = 'file';
    const SyncConsumer  = 'sync';
    const DebugConsumer = 'debug';

    /**
     * @var SensorsAnalytics
     */
    private $instance;


    public function __construct(string $consumerType, string $serverUrl = '', string $projectName = '')
    {
        if (in_array($consumerType, [self::SyncConsumer, self::DebugConsumer]) && !$serverUrl) {
            throw new Exception("sync/debug 模式服务端url不能为空");
        }

        switch ($consumerType) {
            case self::FileConsumer:
                $filePath = storage_path(sprintf("shenCe/%s-track.log", date('Y-m-d')));
                $consumer = new FileConsumer($filePath);
                break;
            case self::SyncConsumer:
                $consumer = new BatchConsumer($serverUrl);
                break;
            case self::DebugConsumer:
                # 从神策分析配置页面中获取的数据接收的 URL
                # 初始化一个 Consumer，用于数据发送
                $consumer = new DebugConsumer($serverUrl, true);
                break;
            default:
                throw new Exception("无效的消费类型");
        }

        if (!$projectName) {
            $this->instance = new SensorsAnalytics($consumer);
        } else {
            $this->instance = new SensorsAnalytics($consumer, $projectName);
        }
    }

    public function instance(): SensorsAnalytics
    {
        return $this->instance;
    }

}