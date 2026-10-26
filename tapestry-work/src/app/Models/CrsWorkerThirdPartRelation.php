<?php
/**
 * Created by PhpStorm.
 * User: Zhujiaqi
 * Date: 2020/9/29
 * Time: 12:02
 */

namespace App\Models;

class CrsWorkerThirdPartRelation extends Base
{
    //同步来源：1-crs, 2-hwj
    const SYNC_SOURCE_CRS = 1;
    const SYNC_SOURCE_HWJ = 2;

    const SOURCE_LIST = [
        self::SYNC_SOURCE_CRS,
        self::SYNC_SOURCE_HWJ
    ];

    protected $table = 'crs_worker_third_part_relation';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'sync_source',
        'enterprise_id',
        'third_part_user_id',
        'store_staff_id',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

    /**
     * @var array
     */
    public static $selectField = [
        'id',
        'sync_source',
        'enterprise_id',
        'third_part_user_id',
        'store_staff_id',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

    public function getStoreStaffIdByThirdPartUserId($thirdPartUserId, $enterpriseId) {
        $queryData =  self::query()
            ->where([
                ['third_part_user_id', $thirdPartUserId],
                ['enterprise_id', $enterpriseId],
            ])
            ->value('store_staff_id');
        return $queryData ?? 0;

    }
}
