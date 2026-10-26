<?php
/**
 * 微信优惠卷
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-24
 * Time: 18:28
 */

namespace App\Http\Controllers\V1\AppCenter\WxCoupon;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxCoupon;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;

class DelController extends BaseController
{
    public function index(Request $request)
    {
        $input                  = $request->only(['id', 'force']);
        $input['enterprise_id'] = $this->enterpriseId;

        $response = (new WxCoupon())->del($input);

        $redisKey = sprintf(Constant::WX_COUPON_DETAIL, intval($input['id']));
        $redis = Redis::connection();
        $redis->del($redisKey);

        return success($response);
    }
}
