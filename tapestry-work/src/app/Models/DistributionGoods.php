<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class DistributionGoods extends Base
{
    /**
     * 推广佣金类型：1=>百分比,2=>金额;
     */
    const PROMOTION_COMMISSION_TYPE_PERCENTAGE = 1;
    const PROMOTION_COMMISSION_TYPE_AMOUNT = 2;

    /**
     * 服务佣金类型：1=>百分比,2=>金额;
     */
    const SERVICE_COMMISSION_TYPE_PERCENTAGE = 1;
    const SERVICE_COMMISSION_TYPE_AMOUNT = 2;

    protected $table = 'distribution_goods';

    protected $primaryKey = 'id';

    /*
     * 品牌上新
     */
    const NEW_GOODS = 1;
}
