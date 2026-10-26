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

/**
 * Class StatisticsController
 *
 * @package App\Http\Controllers\App\V1\Store
 */
class CompassStatisticsController extends Controller
{
    /**
     * 数据罗盘统计
     *
     * @param Request          $request
     * @param StorePerformance $storePerformance
     *
     * @return array|\Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, StorePerformance $storePerformance)
    {
        // 参数校验
        $rules   = [
            'store_id'   => 'required|integer',
            'type'       => 'required|integer|in:1,2',
            'start_time' => 'integer',
            'end_time'   => 'integer',
        ];
        $message = [
            'store_id.required'  => 'store_id参数为必填项',
            'store_id.integer'   => 'store_id参数必须为整型数字',
            'type.required'      => 'type参数为必填项',
            'type.integer'       => 'type参数必须为整型数字',
            'type.in'            => 'type类型不符合规范',
            'start_time.integer' => 'start_time参数必须为整型数字',
            'end_time.integer'   => 'end_time参数必须为整型数字',
        ];
        $this->validate($request, $rules, $message);
        $enterpriseId = Auth::user()['enterprise_id'];
        $storeId      = $request->post('store_id');
        $startTime    = $request->post('start_time', '');
        $endTime      = $request->post('end_time', '');
        $type         = $request->post('type', 1);
        $workerIds    = $request->post('worker_ids', []);

        $returnData = $storePerformance->storeStatistics([
            'enterprise_id' => $enterpriseId,
            'store_id'      => $storeId,
            'worker_ids'    => $workerIds,
            'start_time'    => $startTime,
            'end_time'      => $endTime,
            'type'          => $type
        ]);

        return success($returnData);
    }
}
