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
use Illuminate\Validation\Rule;

class ListController extends BaseController
{

    /**
     * 获取关联列表
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        // 参数验证规则
        $rules   = [
            'status' => ['required', Rule::in([-1, 0, 1])],
        ];
        $message = [
            'status.*' => '关联状态有误',
        ];
        $this->validate($request, $rules, $message);
        $input                  = $request->only(['coupon_stock_id', 'merchant_coupon_stock_id', 'merchant_coupon_name', 'status', 'page', 'size']);
        $input['enterprise_id'] = $this->enterpriseId;
        $data                   = (new CouponRelation())->list($input);
        return success($data);
    }

}
