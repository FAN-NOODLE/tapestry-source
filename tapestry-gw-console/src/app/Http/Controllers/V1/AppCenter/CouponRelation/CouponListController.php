<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\CouponRelation;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CouponRelation;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class CouponListController extends BaseController
{

    /**
     * 获取可关联的代金券列表
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $input                  = $request->only(['coupon_stock_id', 'page', 'size']);
        $input['enterprise_id'] = $this->enterpriseId;
        $data                   = (new CouponRelation())->couponList($input);
        return success($data);
    }

}
