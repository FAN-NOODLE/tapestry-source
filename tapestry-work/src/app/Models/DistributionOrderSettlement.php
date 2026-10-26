<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class DistributionOrderSettlement extends Base
{
    protected $table = 'distribution_order_settlement';

    protected $primaryKey = 'id';

    /**
     * 结算状态：0待结算 1 已结算 2已取消
     *
     * @var int
     */
    const  WAIT_SETTLE = 0;
    const  ADOPT_SETTLE = 1;
    const  CANCEL_SETTLE = 2;

    static public $status = [
        '0'=>'待结算',
        '1'=>'已结算',
        '2'=>'已取消'
    ];

    protected $fillable = [
        'id',
        'enterprise_id',
        'order_no',
        'order_status',
        'settlement_status',
        'settlement_standard',
        'order_pay_time',
        'promotion_worker_id',
        'promotion_commission',
        'service_shopping_guid_id',
        'service_commission',
        'partner_id',
        'partner_commission',
        'partner_enterprise_id',
        'created_at',
        'updated_at',
        'service_commission_snapshot',
        'partner_commission_snapshot',
        'promotion_commission_snapshot'
    ];


}
