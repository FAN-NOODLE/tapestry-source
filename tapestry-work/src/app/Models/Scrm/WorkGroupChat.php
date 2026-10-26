<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: WorkGroupChat.php
 */

namespace App\Models\Scrm;


class WorkGroupChat extends ScrmInit
{
    /**
     * 群状态（1-存在，2-解散）
     */
    const GROUP_EXIST_STATUS = 1;
    const GROUP_DISMISS_STATUS = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'work_group_chat';


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
     * 获取群组数据
     * @param $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $memberIds
     * @param string $order
     * @param array $fields
     * @param int $limit
     * @param bool $page
     * @return array
     * @Date: 2021/6/29
     * @Time: 20:36
     */
    public function getGroup(int $enterpriseId, int $startStatTime = 0, int $endStatTime = 0, array $memberIds = [], string $order = '', array $fields = ['*'], int $limit = 10, bool $page = false) : array
    {
        $memberIds = array_filter(array_unique(array_column($memberIds, 'id')));

        $query = self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->when(!empty($memberIds), function ($query) use ($memberIds) {
                return $query->whereIn('owner_member_id', $memberIds);
            })
            ->when(!empty($startStatTime), function ($query) use ($startStatTime) {
                return $query->where('create_time', '>=', $startStatTime);
            })
            ->when(!empty($endStatTime), function ($query) use ($endStatTime) {
                return $query->where('create_time', '<=', $endStatTime);
            })
            ->where('status', WorkGroupChat::GROUP_EXIST_STATUS)
            ->when(!empty($order), function ($query) use ($order) {
                return $query->orderBy($order, 'DESC');
            });

        if ($page) {
            return $query->paginate($limit)->toArray();
        }

        return $query->get()->toArray();
    }

}