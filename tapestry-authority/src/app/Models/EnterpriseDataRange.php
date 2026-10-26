<?php


namespace App\Models;

class EnterpriseDataRange extends Base
{

    /**
     * -1- 不限 0-全部门店 1-区域 2-门店 3-商城 4-部门
     */
    const TYPE_NO_LIMIT = -1;
    const TYPE_ALL = 0;
    const TYPE_AREA = 1;
    const TYPE_STORE = 2;
    const TYPE_SHOP = 3;
    const TYPE_DEPARTMENT = 4;

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_data_range';

    protected $fillable = [
        'id',
        'enterprise_id',
        'enterprise_admin_id',
        'type',
        'object_id',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    /**
     * 获取数据权限
     * @param $adminId
     * @param $oId
     * @param $type
     * @return EnterpriseDataRange|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public static function getDataRange($adminId, $oId, $type)
    {
        return self::query()
            ->where('enterprise_admin_id', '=', $adminId)
            ->where('object_id', '=', $oId)
            ->where('type', '=', $type)
            ->first();
    }

    public function accountDataRange($adminId)
    {
        return self::query()
            ->select(['id', 'enterprise_id', 'type', 'object_id'])
            ->where('enterprise_admin_id', $adminId)
            ->get()
            ->toArray();
    }

    public function accountDataRangeByEid($adminId)
    {
        $data  = self::query()
            ->select(['id', 'enterprise_id', 'type', 'object_id'])
            ->where('enterprise_id', $adminId)
            ->get()
            ->toArray();
        $array = [];
        $list  = [];
        foreach ($data as $key => $value) {
            if (in_array($value['object_id'], $array)) {
                continue;
            }
            $array[] = $value['object_id'];
            $list[]  = $value;
        }
        return $list;
    }
}