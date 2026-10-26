<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: WorkGroupChatMember.php
 */

namespace App\Models\Scrm;

use App\Handlers\SoftDeletingScope;
use Illuminate\Support\Facades\DB;

class WorkGroupChatMember extends ScrmInit
{
    /**
     * 成员类型（企业成员-1, 外部联系人-2）
     */
    const ENTERPRISE_MEMBERS_TYPE = 1;
    const EXTERNAL_CONTACTS_TYPE = 2;

    /**
     * 群状态（入群-1，退群-2）
     */
    const JOIN_GROUP = 1;
    const LEAVE_GROUP = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'work_group_chat_member';


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
     * 获取所有在群的群组成员数据
     * @param $enterpriseId
     * @param $startStatTime
     * @param $endStatTime
     * @param $groupChatId
     * @param array $fields
     * @param string $group
     * @param int $status
     * @return array
     * @Date: 2021/6/29
     * @Time: 20:37
     */
    public function getGroupMember(int $enterpriseId, int $startStatTime, int $endStatTime, array $groupChatId, array $fields = ['*'], string $group = '', int $status = 0) : array
    {
        return self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('group_chat_id', $groupChatId)
            ->when(!empty($startStatTime), function ($query) use ($startStatTime) {
                return $query->where('join_time', '>=', $startStatTime);
            })
            ->when(!empty($endStatTime), function ($query) use ($endStatTime) {
                return $query->where('join_time', '<=', $endStatTime);
            })
            ->when(!empty($status), function ($query) use ($status) {
                return $query->where('status', $status);
            })
            ->where('type', self::EXTERNAL_CONTACTS_TYPE)
            ->when(!empty($group), function ($query) use ($group) {
                return $query->groupBy($group);
            })
            ->get()
            ->toArray();
    }

    /**
     * 获取最近入群与离群的群组成员
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $groupChatId
     * @param int $status
     * @param array $fields
     * @param string $group
     * @return array
     * @Date: 2021/6/29
     * @Time: 20:38
     */
    public function getJoinOrLeaveGroup(int $enterpriseId, int $startStatTime, int $endStatTime, array $groupChatId, int $status, array $fields = ['*'], string $group = '') : array
    {
        return self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('group_chat_id', $groupChatId)
            ->where('type', self::EXTERNAL_CONTACTS_TYPE)
            ->when(!empty($status), function ($query) use ($status, $startStatTime, $endStatTime) {
                $fields = $status == self::JOIN_GROUP ? 'join_time' : 'leave_time';

                return $query->where('status', $status)
                    ->where($fields, '>=', $startStatTime)
                    ->where($fields, '<=', $endStatTime);
            })
            ->when(!empty($group), function ($query) use ($group) {
                return $query->groupBy($group);
            })
            ->get()
            ->toArray();
    }

    public function getGroupMemberCount(int $enterpriseId, int $startStatTime, int $endStatTime, array $groupChatId, array $fields = ['*'], string $group = '', int $status = 0) : array
    {
        $objQuery = self::query()
            ->withoutGlobalScope(new SoftDeletingScope())
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('group_chat_id', $groupChatId)
            ->where('deleted_at', 0)
            ->when(!empty($startStatTime), function ($query) use ($startStatTime) {
                return $query->where('join_time', '>=', $startStatTime);
            })
            ->when(!empty($endStatTime), function ($query) use ($endStatTime) {
                return $query->where('join_time', '<=', $endStatTime);
            })
            ->when(!empty($status), function ($query) use ($status) {
                return $query->where('status', $status);
            })
            ->where('type', self::EXTERNAL_CONTACTS_TYPE)
            ->when(!empty($group), function ($query) use ($group) {
                return $query->groupBy($group);
            });

        $result = DB::table(DB::raw("({$objQuery->toSql()}) as c"))
            ->select(DB::raw("sum( c.groupMemberCount ) as groupMemberCount"))
            ->mergeBindings($objQuery->getQuery())
            ->get()
            ->toArray();

        return $result;
    }

}