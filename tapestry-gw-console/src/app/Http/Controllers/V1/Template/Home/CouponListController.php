<?php

namespace App\Http\Controllers\V1\Template\Home;

use CrsSdk\Shop\Template;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class CouponListController extends BaseController
{
    /**
     * 优惠券弹窗 - 包含微信和店内
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->all();

        $params = array_merge($params, ['enterprise_id' => $this->enterpriseId]);

        $list = (new Template())->couponList($params);

        return success($list);
    }
}
