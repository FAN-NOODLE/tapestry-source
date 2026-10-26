<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 18:36
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\StorePerformance;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

/**
 * 员工数据统计
 * Class StatisticsController
 * @package App\Http\Controllers\App\V1\Store
 */
class StatisticsController extends Controller
{
    /**
     * 员工业绩数据统计
     * @param Request $request
     * @param StorePerformance $storePerformance
     * @return array|\Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, StorePerformance $storePerformance)
    {
        // 参数校验
        $rules = [
            'store_id'      => 'required|integer',
            'worker_id'     => 'required|integer',
            'start_time'    => 'integer',
            'end_time'      => 'integer',
        ];
        $message = [
            'store_id.required'      => 'store_id参数为必填项',
            'store_id.integer'       => 'store_id参数必须为整型数字',
            'worker_id.required'     => 'worker_id参数不能为空',
            'worker_id.integer'      => 'worker_id参数必须为整型数字',
            'start_time.integer'     => 'start_time参数必须为整型数字',
            'end_time.integer'       => 'end_time参数必须为整型数字',
        ];
        $this->validate($request, $rules, $message);
        $enterpriseId = Auth::user()['enterprise_id'];
        $storeId = $request->get('store_id');
        $workerId = $request->get('worker_id');
        $startTime = $request->get('start_time', time());
        $endTime = $request->get('end_time', time());
        $returnData = [
            'in_sales_order_num'   => 0,
            'complete_order_num'   => 0,
            'sale_all_num'         => 0,
            'in_sales_money_total' => 0,
            'complete_money_total' => 0,
            'sale_all_money_total' => 0,
            'real_money_total'     => 0,
        ];
        try {
            $storePerformanceData = $storePerformance->statistics([
                'enterprise_id' => $enterpriseId,
                'store_id'      => $storeId,
                'worker_id'     => $workerId,
                'start_time'    => date('Ymd', $startTime),
                'end_time'      => date('Ymd', $endTime),
            ]);
            $storePerformanceData = array_column($storePerformanceData, null, 'worker_id');
            $returnData = $storePerformanceData[$workerId] ?? $returnData;
            return success($returnData);
        } catch (\Exception $e) {
            error(655004);
        }
    }
}
