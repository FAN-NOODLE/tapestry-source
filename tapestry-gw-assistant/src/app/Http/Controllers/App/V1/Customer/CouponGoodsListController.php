<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/17
 * Time: 17:59
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\Coupon;
use Illuminate\Http\Request;

class CouponGoodsListController extends BaseController {

    /**
     * 顾客优惠券详情
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params                 = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $result                 = (new Coupon())->goodsList($params);
        return success($result);
    }

}