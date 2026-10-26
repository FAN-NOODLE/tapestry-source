<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-15
 * Time: 12:55
 */

namespace App\Models;

class StoreMemberViewTab extends BaseModel
{

    /**
     * 查看类型 1 抢服务订单 2 服务订单待接单
     */
    const EXTENSION_TAB_TYPE = 1;
    const SERVICE_ORDER_TO_BE_RECEIVED_TAB_TYPE = 2;
    const CUSTOMER_TAB_TYPE = 3;
    const CUSTOMER_TOTAL_TAB_TYPE = 4;

    protected $table = 'store_member_view_tab';

    protected $fillable = [
        'id',
        'enterprise_id',
        'store_id',
        'store_worker_id',
        'extension_view_time',
        'type'
    ];

    /**
     * 获取浏览时间
     * @param $params
     * @return StoreMemberViewTab|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function getViewTime($params)
    {
        $detail = self::query()
            ->where($params)
            ->first();
        return $detail;
    }
}
