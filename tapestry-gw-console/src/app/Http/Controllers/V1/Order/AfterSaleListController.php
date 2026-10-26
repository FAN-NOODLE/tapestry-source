<?php
/**
 * Created by PhpStorm.
 * User: zq
 * Date: 19-4-8
 * Time: 下午6:20
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Order;

class AfterSaleListController extends BaseController
{
    /**
     * 售后列表
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input = $request->only(['page', 'size', 'search']);
        $input['enterprise_id'] = $this->enterpriseId;

        // 超管（1=== $this->hasShopAuth）可查看企业下所有订单，否则，只可查看当前用户权限内的店铺的订单
        if (2 === $this->hasShopAuth) {
            $input['search']['store_ids'] = $this->storeIds ?: [];
        }

        $res = (new Order())->afterSaleList($input);
        return success($res);
    }
}
