<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\Coupon;

use App\Exceptions\Error;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\CouponRequest;
use CrsSdk\Shop\Coupon;

class AddController extends BaseController
{
    /**
     * 礼劵验证器
     *
     * @var CouponRequest
     */
    protected $couponRequest;

    /**
     * AddController constructor.
     * @param CouponRequest $couponRequest
     */
    public function __construct(CouponRequest $couponRequest)
    {
        parent::__construct();

        // 礼劵验证器赋值
        $this->couponRequest = $couponRequest;
    }

    /**
     * 礼劵添加
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/30 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['name', 'discount_amount', 'condition_value', 'start_time', 'end_time', 'range', 'channel']);

        // 请求数据验证
        $validator = $this->couponRequest->common($params);
        if ($validator->fails()) {
            throw new Error(400000, $validator->errors()->first());
        }

        $params['enterprise_id'] = $this->enterpriseId;
        if (isset($params['range'])) {
            $params['range'] = !empty($params['range']) ? json_encode($params['range']) : '';
        }
        $params['type']      = 1;
        $params['use_type']  = 1;
        $params['effective'] = 2;
        $params['limit']     = 1;
        try {
            $this->writeOperatorLog('新增礼券:'.$params['name']);

            (new Coupon())->add($params);
            return success([]);
        } catch (\Exception $e) {
            return error(710001);
        }
    }
}
