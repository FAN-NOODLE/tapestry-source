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

class ListController extends BaseController
{
    /**
     * 礼劵列表
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/25 11:16
     */
    public function index(Request $request)
    {
        $params = $request->input();

        // 查询参数处理
        isset($params['name']) && $where['name']       = $params['name'];
        isset($params['status']) && $where['status']   = $params['status'];
        isset($params['channel']) && $where['channel'] = $params['channel'];

        $where['page']          = isset($params['page']) ? $params['page'] : 1;
        $where['limit']         = isset($params['limit']) ? $params['limit'] : 10;
        $where['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Coupon())->lists($where);
            return success($list);
        } catch (\Exception $e) {
            error(710004);
        }
    }
}
