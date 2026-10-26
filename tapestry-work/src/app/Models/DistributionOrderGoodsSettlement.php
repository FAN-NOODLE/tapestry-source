<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class DistributionOrderGoodsSettlement extends Base
{
    protected $table = 'distribution_order_goods_settlement';

    protected $primaryKey = 'id';

    /**
     * 订单商品退款状态  refund_status 默认 0 部分 1 全部 2
     */
    const  REFUND_STATUS_M = 0;
    const  REFUND_STATUS_PART = 1;
    const  REFUND_STATUS_ALL = 2;
}
