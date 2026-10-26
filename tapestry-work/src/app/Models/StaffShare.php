<?php


namespace App\Models;


class StaffShare extends Base
{
    /**
     * 绑定
     */
    const BINDING_TYPE_YES = 1;

    /**
     * 解除绑定
     */
    const BINDING_TYPE_NO = 2;

    const BINDING_TYPE_PARENT_UNBINDING = 3;

    const BINDING_TYPE_ONE = 1;
    const BINDING_TYPE_TWO = 2;
    const CHANGE_OWNER_ONE = 1;
    const CHANGE_OWNER_TWO = 2;

    /**
     * 业绩归属 1 我 2 他人 3 无归属
     */
    const ACHIEVEMENT_ASCRIPTION_SELF = 1;
    const ACHIEVEMENT_ASCRIPTION_OTHER = 2;
    const ACHIEVEMENT_ASCRIPTION_NONE = 3;

    /**
     * 员工内购分享表
     *
     * @var string
     */
    protected $table = 'staff_share';

    protected $fillable = [
        'id',
        'enterprise_id',
        'store_id',
        'user_id',
        'type',
        'parent_user_id',
        'mem_uid',
        'worker_id',
        'level',
        'qr_code',
        'binding_type',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

    public function getStaffShareInfoByCondition($where, $bindStatus = true)
    {
        if ($bindStatus) {
            $where = array_merge($where,['binding_type' => self::BINDING_TYPE_YES]);
        }

        $shareInfo = self::query()
            ->where($where)
            ->orderByDesc('created_at')
            ->first();

        return $shareInfo ? $shareInfo->toArray() : [];
    }

}
