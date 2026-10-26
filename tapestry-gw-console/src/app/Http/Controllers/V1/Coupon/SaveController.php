<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\Coupon;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Coupon;

class SaveController extends BaseController
{
    /**
     * 礼劵添加
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/30 11:16
     */
    public function index(Request $request)
    {
        $params = $request->only(['coupon_id', 'name', 'discount_amount', 'condition_value', 'start_time', 'end_time', 'range', 'channel']);

        // 礼劵ID不能为空
        if (!isset($params['coupon_id']) || empty($params['coupon_id'])) {
            error(400001);
        }

        //特殊数据处理
        if (isset($params['range'])) {
            $params['range'] = !empty($params['range']) ? json_encode($params['range']) : '';
        }

        try {
            $this->writeOperatorLog('编辑礼券:'.$params['name']);

            (new Coupon())->save($params);
            return success([]);
        } catch (\Exception $e) {
            error(710005);
        }
    }
}
