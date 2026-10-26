<?php
/**
 * Coach的CIIE吊牌预约活动 (导购端)
 */


namespace App\Http\Controllers\App\V1\Activity;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Work\Activity\ReserveCIIE;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ReserveCIIECheckController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * sa导购白名单
     * @param Request $request
     * @return JsonResponse
     */
    public function WhiteLists(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = (new ReserveCIIE())->WhiteLists($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 核销券
     * @return void
     */
    public function Index(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = (new ReserveCIIE())->WriteOffCoupon($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 获取全部预约列表
     * @return void
     */
    public function GetLists(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = (new ReserveCIIE())->ReserveLists($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 获取预约时间端的配置列表
     * @return void
     */
    public function GetReserveTimeList(Request $request): JsonResponse
    {
        $params = $request->all();
        try {
            $result = (new ReserveCIIE())->TimeLists($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function ReserveTimeConfig(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = (new ReserveCIIE())->TimeConfig($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 导购已核销的列表
     * @param Request $request
     * @return JsonResponse
     */
    public function CheckLists(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = (new ReserveCIIE())->CheckLists($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }


}