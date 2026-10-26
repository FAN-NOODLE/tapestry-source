<?php

namespace App\Http\Controllers\App\V1\Activity;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Work\Activity\Scratch;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;


class CheckCouponController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 核销券
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = (new Scratch())->WriteOffCoupon($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 导购已核销券列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/15 16:26
     */
    public function CouponList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = (new Scratch())->CouponList($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }
}