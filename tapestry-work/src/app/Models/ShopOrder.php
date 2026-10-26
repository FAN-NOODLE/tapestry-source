<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/2\App\Models\UserAction::self::ACTION_TYPE_ORDER
 * Time: 17:45
 */

namespace App\Models;

use App\Common\Constant;
use Illuminate\Database\Eloquent\Relations\HasOne;
use Illuminate\Support\Facades\DB;

class ShopOrder extends Base
{

    // 未成交
    const STATUS_WAIT_ARRAY = [21, 20, 23];

    // 已成交
    const STATUS_PAY_ARRAY = [25, 30, 40, 41, 60, 61, 80, 81, 82, 83, 84, 100, 200, 201];


    protected $table = 'shop_order';

    public function goods()
    {
        return $this->hasMany(ShopOrderGoods::class, 'order_no', 'order_no');
    }

    public function traces()
    {
        return $this->hasMany(ShopOrderTrace::class, 'order_id', 'id');
    }
}
