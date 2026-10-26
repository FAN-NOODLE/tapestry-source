<?php


namespace App\Models;


class StaffRight extends Base
{
    /**
     * 允许内购的对象类型：1-部门 2-员工
     */
    const TYPE_DEP = 1;
    const TYPE_STAFF = 2;

    protected $table = 'staff_right';

    protected $guarded = [];

    /**
     * 根据员工memuid获取数据
     * @param $enterpriseId
     * @param $memUid
     * @return array
     */
    public function getStaffRightFirst($enterpriseId, $memUid)
    {
        $detail = self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('object_extra_id', 'like', '%' . $memUid . '%')
            ->orderByDesc('id')
            ->first(['id']);
        return !empty($detail) ? $detail->toArray() : [];
    }
}
