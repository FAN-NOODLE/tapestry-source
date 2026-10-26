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
use Illuminate\Validation\ValidationException;

class RelationController extends BaseController
{

    /**
     * 关联代金券和商家券
     * @param Request $request
     * @return JsonResponse
     * @throws ValidationException
     * @throws Exception
     */
    public function index(Request $request)
    {
        // 参数验证规则
        $rules   = [
            'coupon_stock_id'          => 'required',
            'merchant_coupon_stock_id' => 'required',
            'status'                   => ['required', Rule::in([0, 1])],
        ];
        $message = [
            'coupon_stock_id.*'          => '微信优惠券批次号有误',
            'merchant_coupon_stock_id.*' => '商家券批次号有误',
            'status.*'                   => '关联状态有误',
        ];
        $this->validate($request, $rules, $message);
        $input                  = $request->only(['coupon_stock_id', 'merchant_coupon_stock_id', 'status']);
        $input['enterprise_id'] = $this->enterpriseId;
        $data                   = (new CouponRelation())->relation($input);
        return success($data);
    }

}
