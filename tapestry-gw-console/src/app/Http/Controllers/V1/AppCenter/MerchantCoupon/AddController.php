<?php
/**
 * 微信优惠卷新增
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-23
 * Time: 17:48
 */

namespace App\Http\Controllers\V1\AppCenter\MerchantCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MerchantCoupon;

class AddController extends BaseController
{

    public function index()
    {
        $input = $this->request->only([
            'coupon_name',
            'coupon_stock_type',
            'start_time',
            'end_time',
            'goods_name',
            'coupon_use_rule',
            'stock_send_rule',
            'display_pattern_info',
            'custom_entrance',
            'merchant_list'
        ]);
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new MerchantCoupon())->add($input);

        return success($list);
    }
}
