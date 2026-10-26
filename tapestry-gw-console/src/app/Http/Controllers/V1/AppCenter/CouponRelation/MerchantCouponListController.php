<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\CouponRelation;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CouponRelation;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class MerchantCouponListController extends BaseController
{

    /**
     * 获取可关联的商家券列表
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $input                  = $request->only(['keywords', 'page', 'size']);
        $input['enterprise_id'] = $this->enterpriseId;
        $data                   = (new CouponRelation())->merchantCouponList($input);
        return success($data);
    }

}
