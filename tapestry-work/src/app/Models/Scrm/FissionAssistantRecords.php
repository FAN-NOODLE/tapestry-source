<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: FissionAssistantRecords.php
 */

namespace App\Models\Scrm;


class FissionAssistantRecords extends ScrmInit
{
    /**
     * 入群状态：0-未入群，1-已入群，2-已退群
     */
    const NOT_JOIN_CHAT_GROUP = 0;
    const ALREADY_JOIN_CHAT_GROUP = 1;
    const LEAVE_CHAT_GROUP = 2;

    /**
     * 添加好友记录：0-未添加，1-已添加，2-已删除
     */
    const NOT_ADD_MEMBER = 0;
    const ALREADY_ADD_MEMBER = 1;
    const DELETE_ADD_MEMBER = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'fission_assistant_records';


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
     * 获取参加的活动
     * @param $enterpriseId
     * @param $startStatTime
     * @param $endStatTime
     * @param $groupActivityIds
     * @param array $fields
     * @param string $order
     * @param string $group
     * @return array
     * @Date: 2021/6/29
     * @Time: 20:33
     */
    public function getJoinGroupActivity($enterpriseId, $startStatTime, $endStatTime, $groupActivityIds, $fields = ['*'], $order = '', $group = '')
    {
        return self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('activity_id', $groupActivityIds)
            ->where('join_chat_group_status', self::ALREADY_JOIN_CHAT_GROUP)
            ->where('deleted_at', 0)
            ->when(!empty($order), function ($query) use ($order) {
                return $query->orderBy($order, 'DESC');
            })
            ->when(!empty($group), function ($query) use ($group) {
                return $query->groupBy($group);
            })
            ->get()
            ->toArray();
    }
}