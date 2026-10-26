<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 18:36
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\StorePerformance;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use mysql_xdevapi\Exception;

/**
 * 门店数据统计
 * Class StatisticsController
 * @package App\Http\Controllers\App\V1\Store
 */
class StatisticsController extends Controller
{
    /**
     * 门店业绩数据统计
     * @param Request $request
     * @param StorePerformance $storePerformance
     * @return array|\Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, StorePerformance $storePerformance)
    {
        // 参数校验
        $rules = [
            'store_id' => 'required|integer',
            'start_time' => 'integer',
            'end_time' => 'integer',
        ];
        $message = [
            'store_id.required' => 'store_id参数为必填项',
            'store_id.integer' => 'store_id参数必须为整型数字',
            'start_time.integer' => 'start_time参数必须为整型数字',
            'end_time.integer' => 'end_time参数必须为整型数字',
        ];
        $input = $this->validate($request, $rules, $message);
        $enterpriseId = Auth::user()['enterprise_id'];
        $storeId = $request->get('store_id');
        $startTime = $request->get('start_time', time());
        $endTime = $request->get('end_time', time());
        $default = [
            'in_sales_order_num' => 0,
            'complete_order_num' => 0,
            'sale_all_num' => 0,
            'in_sales_money_total' => 0,
            'complete_money_total' => 0,
            'sale_all_money_total' => 0,
            'real_money_total' => 0,
            'store_id' => $storeId
        ];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);
        try {
            $storePerformanceData = $storePerformance->statistics([
                'enterprise_id' => $enterpriseId,
                'store_id'      => $storeId,
                'start_time'    => date('Ymd', $startTime),
                'end_time'      => date('Ymd', $endTime),
                'worker_id'     => $workerDetail['id']
            ]);
            $storePerformanceData = array_column($storePerformanceData, null, 'store_id');
            $default = $storePerformanceData[$storeId] ?? $default;
            return success($default);
        } catch (\Exception $e){
            error(653002);
        }
    }
}
