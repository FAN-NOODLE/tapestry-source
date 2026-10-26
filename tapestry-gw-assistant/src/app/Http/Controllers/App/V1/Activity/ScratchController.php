<?php

namespace App\Http\Controllers\App\V1\Activity;

use App\Http\Controllers\Controller;
use App\Micros\Work\Activity\Scratch;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

/**
 * 刮刮乐活动
 */
class ScratchController extends Controller
{

    /**
     * @var Scratch
     */
    private $scratch;

    public function __construct()
    {
        parent::__construct();
        $this->scratch = new Scratch();
    }

    /**
     * 活动详情
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:55
     */
    public function Info(Request $request): JsonResponse
    {
        $tokenData               = $request->get('tokenData');
        $params                  = $request->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];
        try {
            $result = $this->scratch->Info($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        $this->scratch->UserView($params);
        return success($result);
    }


    /**
     * 刮奖
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:55
     */
    public function ScratchPrize(Request $request): JsonResponse
    {
        $tokenData               = $request->get('tokenData');
        $params                  = request()->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];
        $params['unionid']       = $tokenData['union_id'];
        $result                  = $this->scratch->ScratchPrize($params);

        return success($result);
    }

    /**
     * 券列表
     *
     * @return JsonResponse
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:44
     */
    public function CouponList(): JsonResponse
    {
        $tokenData               = $this->request->get('tokenData');
        $params                  = request()->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];

        $result = $this->scratch->CouponList($params);

        return success($result);
    }

    /**
     * 检查用户券核销状态
     * @return JsonResponse
     * @throws Exception
     */
    public function CheckCoupon(): JsonResponse
    {
        $tokenData               = $this->request->get('tokenData');
        $params                  = request()->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];

        $result = $this->scratch->CheckCoupon($params);

        return success($result);
    }


    /**
     * @return JsonResponse
     * @throws Exception
     */
    public function ShareInfo(): JsonResponse
    {
        $tokenData               = $this->request->get('tokenData');
        $params                  = request()->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];

        $result = $this->scratch->ShareInfo($params);
        return success($result);
    }

}