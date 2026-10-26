<?php

namespace App\Events;

class StoreOrderStatus extends Event
{
    /**
     * 队列传递过来的消息数据
     *
     * @var array
     */
    private $queueData = [];

    /**
     * 任务开始时间
     *
     * @var float
     */
    private $jobStartTime = 0;

    /**
     * StoreOrderStatus constructor.
     *
     * @param array $queueData
     * @param float $jobStartTime
     */
    public function __construct(array $queueData, float $jobStartTime)
    {
        $this->queueData = $queueData;
        $this->jobStartTime = $jobStartTime;
    }

    /**
     * 获取任务开始执行时间
     *
     * @return float|int
     */
    public function getJobStartTime()
    {
        return $this->jobStartTime;
    }

    /**
     * 获取队列数据
     *
     * @return array
     */
    public function getQueueData()
    {
        return $this->queueData;
    }
}
