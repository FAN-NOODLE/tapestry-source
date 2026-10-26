<?php

namespace App\Http\Controllers\Html\V1\MarketingPromotion;

use App\Http\Controllers\Html\V1\BaseController;
use CrsSdk\Shop\WxCoupon;

class CouponListController extends BaseController
{
    /**
     * 获取优惠券列表
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['page', 'size']);

        $search = [
            'enterprise_id' => $this->enterpriseId,
            'page' => $params['page'] ?? 1,
            'size' => $params['size'] ?? 10,
            'search' => [
                'time_overdue' => time(),
            ]
        ];

        switchAppHostByDomain($this->domain);
        
        $list = (new WxCoupon())->wxCouponList($search);

        return success([
            'data' => $list['list'] ?? [],
            'total' => $list['total'] ?? 0,
            'page' => $search['page'],
            'size' => $search['size']
        ]);
    }
}
