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

class UpdateController extends BaseController
{
    /**
     * 活动编辑
     * @author 胡龙阳
     * @date 2020/5/28 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->only(['id', 'name', 'activity_type', 'exclude_type', 'exclude_activity_name', 'activity_goods_type', 'start_time', 'end_time', 'buy_limit_num', 'up_status', 'up_type', 'up_time', 'activity_status', 'number_price_type', 'rule', 'goods_list', 'use_rule_list', 'send_gift_stock_list']);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['goods_list'] = $params['goods_list'] ?? [];

        (new DiscountActivity())->update($params);

        return success();
    }

    /**
     * 修改库存
     *
     * @param Request $request
     * @return void
     * @author narutor <550655976@qq.com>
     */
    public function updateFullGiftUpStock(Request $request)
    {
        // 查询参数处理
        $params = $request->only(['id', 'sku_id', 'type', 'num']);

        $return = (new DiscountActivity())->updateFullGiftUpStock($params);

        return success($return);
    }

    /**
     * 修改礼品状态
     *
     * @param Request $request
     * @return void
     * @author narutor <550655976@qq.com>
     */
    public function updateGiftStatus(Request $request)
    {
        // 查询参数处理
        $params = $request->only(['id', 'status']);

        (new DiscountActivity())->updateGiftStatus($params);

        return success();
    }
}
