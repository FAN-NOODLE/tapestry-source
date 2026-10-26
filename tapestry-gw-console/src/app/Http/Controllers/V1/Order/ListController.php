<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use App\Common\Constant;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use Exception;

class ListController extends BaseController
{
    /**
     * 订单列表
     *
     * @param Request $request
     *
     * @return mixed
     * @throws Exception
     */
    public function index(Request $request)
    {
        $input = $request->only(['search', 'page', 'size']);
        $input['search']['enterprise_id'] = $this->enterpriseId; // 企业ID
        $input['search']['domain'] = $this->qyDomain; // 企业domain
        $input['search']['identifier'] = $this->appIdentifier; // 应用标识
        $input['search']['is_goods'] = Constant::ORDER_GOODS_FALSE; // 订单是否包含商品信息：否
        $input['store'] = $this->getStores($this->enterpriseAdmin['role'], false) ?? [];

        // 超管(1 === $this->hasShopAuth)可查看企业下所有订单，否则，只可查看当前用户权限内的店铺的订单
        if (2 === $this->hasShopAuth) {
            $input['search']['store_ids'] = $this->storeIds ?: [];
        }

        $orderList = (new Order())->orderList($input);
        
        return success($orderList);
    }
}
