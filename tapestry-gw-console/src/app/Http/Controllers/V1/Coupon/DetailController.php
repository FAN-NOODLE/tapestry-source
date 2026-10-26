<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\Coupon;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Coupon;

class DetailController extends BaseController
{
    /**
     * 礼劵详情
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/30 11:16
     */
    public function index(Request $request)
    {
        $params = $request->only(['coupon_id']);

        // 礼劵ID不能为空
        if (!isset($params['coupon_id']) || empty($params['coupon_id'])) {
            error(400001);
        }

        try {
            $detail = (new Coupon())->detail(['coupon_id' => $params['coupon_id']]);
            return success($detail);
        } catch (\Exception $e) {
            error(710003);
        }
    }
}
