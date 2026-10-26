<?php
/**
 * 微信优惠卷启用禁用
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-23
 * Time: 17:48
 */
namespace App\Http\Controllers\V1\AppCenter\WxCoupon;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxCoupon;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;

class ChangeStatusController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->only(['id', 'coupon_status']);
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new WxCoupon())->status($input);

        $redisKey = sprintf(Constant::WX_COUPON_DETAIL, intval($input['id']));
        $redis = Redis::connection();
        $redis->del($redisKey);

        return success($list);
    }
}