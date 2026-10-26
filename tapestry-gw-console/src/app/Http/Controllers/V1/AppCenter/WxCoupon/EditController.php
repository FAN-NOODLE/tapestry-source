<?php
/**
 * 微信优惠卷编辑
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-23
 * Time: 17:48
 */

namespace App\Http\Controllers\V1\AppCenter\WxCoupon;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxCoupon;
use Illuminate\Support\Facades\Redis;

class EditController extends BaseController
{

    /**
     * 是否可以叠加使用  1:可以叠加
     */
    const IS_OVERLAY_TRUE = 1;

    /**
     * 是否可以叠加使用 2:不能叠加
     */
    const IS_OVERLAY_FALSE = 2;

    /**
     * 优惠卷编辑
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'coupon_stock_id',
            'is_overlay',
            'overlay_coupons',
            'range',
            'is_all',
            'range_hide',
            'display_home',
            'tag_name',
            'tag_goods',
            'use_day_limit',
            'tag_rule_type',
            'remark'
        ]);

        if (isset($input['use_day_limit']) === false) {
            $input['use_day_limit'] = 0;
        }

        if (! isset($input['use_day_limit'])) {
            error(770002);
        }

        if ((int)$input['use_day_limit'] < 0) {
            error(770003);
        }

        $input['enterprise_id'] = $this->enterpriseId;

        if (isset($input['tag_name'])) {
            if (strlen($input['tag_name']) > 32) {
                error(770001);
            }
        }

        $list = (new WxCoupon())->edit($input);

        $redisKey = sprintf(Constant::WX_COUPON_DETAIL, intval($input['id']));
        $redis = Redis::connection();
        $redis->del($redisKey);

        return success($list);
    }
}
