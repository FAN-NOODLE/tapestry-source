<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: FissionChatGroup.php
 */

namespace App\Models\Scrm;


class FissionChatGroup extends ScrmInit
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'fission_chat_group';


    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $guarded = [];

    /**
     * 获取群组相关的活动
     * @param $enterpriseId
     * @param $startStatTime
     * @param $endStatTime
     * @param $groupChatId
     * @return array
     * @Date: 2021/6/29
     * @Time: 20:35
     */
    public function getGroupActivity(int $enterpriseId, int $startStatTime, int $endStatTime, array $groupChatId) : array
    {
        return self::query()
            ->select(['fission_chat_group.activity_id', 'fission_chat_group.chat_id'])
            ->join('fission_activity', 'fission_chat_group.activity_id', '=', 'fission_activity.id')
            ->where('fission_chat_group.enterprise_id', $enterpriseId)
            ->whereIn('fission_chat_group.chat_id', $groupChatId)
            ->when(!empty($startStatTime) && !empty($endStatTime), function ($query) use ($startStatTime, $endStatTime) {
                return $query->where('fission_activity.status', FissionActivity::FISSION_ACTIVITY_ON)
                    ->where([
                        ['fission_activity.start_time', "<=", $endStatTime],
                        ['fission_activity.end_time', ">=", $startStatTime]
                    ]);
            })
            ->get()
            ->toArray();
    }
}