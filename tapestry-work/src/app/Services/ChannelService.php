<?php

namespace App\Services;

use App\Models\Channel;

class ChannelService extends BaseService
{
    /**
     * 昨天
     */
    const TYPE_YESTERDAY = 1;

    /**
     * 近7天
     */
    const TYPE_NEARLY_7_DAYS = 2;

    /**
     * 近30天
     */
    const TYPE_NEARLY_30_DAYS = 3;

    /**
     * @var Channel
     */
    protected $channel;

    public function __construct(Channel $channel)
    {
        $this->channel = $channel;
    }

    /**
     * 获取查询的起止时间
     *
     * @param array $input
     *
     * @return array
     */
    public function getTimeRang(array $input)
    {
        // 开始时间，默认今天
        $start = strtotime(date('Y-m-d 0:0:0', time()));
        // 当前时间
        $end = time();
        if (!empty($input['start_time']) && !empty($input['end_time'])) {
            $start = $input['start_time'];
            $end   = $input['end_time'];
        }

        if (!empty($input['type'])) {
            $type = intval($input['type']);
            switch ($type) {
                case self::TYPE_YESTERDAY:
                    $start = strtotime(date('Y-m-d 0:0:0', strtotime('-1 day')));
                    $end   = strtotime(date('Y-m-d 23:59:59', strtotime('-1 day')));
                    break;
                case self::TYPE_NEARLY_7_DAYS:
                    $start = strtotime(date('Y-m-d 0:0:0', strtotime('-6 day')));
                    break;
                case self::TYPE_NEARLY_30_DAYS:
                    $start = strtotime(date('Y-m-d 0:0:0', strtotime('-29 day')));
                    break;
                default:
                    break;
            }
        }

        return [$start, $end];
    }

    /**
     * 获取指定日期段内每一天的日期
     * @param int $startTime 开始时间
     * @param int $endTime 结束时间
     * @return array
     */
    function getDateRange($startTime, $endTime)
    {
        $date = [];
        while ($startTime <= $endTime) {
            $date[]    = date('Y-m-d', $startTime);
            $startTime = $startTime + 86400;
        }
        return $date;
    }

    public static function addChannel(array $data)
    {
        $existChannel = Channel::query()
            ->where([
                'enterprise_id' => $data['enterprise_id'],
                'channel'       => $data['channel']
            ])->first(['id']);

        if ($existChannel) {
            return $existChannel->id;
        }
        $ret = Channel::query()->create($data);
        return $ret->id;
    }
}
