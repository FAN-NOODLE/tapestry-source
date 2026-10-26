<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\DiscountActivity;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\DiscountActivity;

class CreateController extends BaseController
{
    /**
     * 活动创建
     * @author 胡龙阳
     * @date 2020/5/28 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['name', 'activity_type', 'exclude_type', 'exclude_activity_name', 'activity_goods_type', 'start_time', 'end_time', 'buy_limit_num', 'up_status', 'up_type', 'up_time', 'activity_status', 'number_price_type', 'rule', 'goods_list', 'use_rule_list', 'send_gift_stock_list']);
        $params['enterprise_id'] = $this->enterpriseId;
        
        $params['goods_list'] = $params['goods_list'] ?? [];

        (new DiscountActivity())->create($params);

        return success();
    }
}
