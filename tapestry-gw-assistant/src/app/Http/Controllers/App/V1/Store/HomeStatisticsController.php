<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 18:36
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\StorePerformance;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

/**
 * Class StatisticsController
 *
 * @package App\Http\Controllers\App\V1\Store
 */
class HomeStatisticsController extends Controller
{
    /**
     * 员工业绩数据统计
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
            'type'       => 'required|integer|in:1,2',
            'start_time' => 'integer',
            'end_time'   => 'integer',
        ];
        $message = [
            'type.required'      => 'type参数为必填项',
            'type.integer'       => 'type参数必须为整型数字',
            'type.in'            => 'type类型不符合规范',
            'start_time.integer' => 'start_time参数必须为整型数字',
            'end_time.integer'   => 'end_time参数必须为整型数字',
        ];
        $this->validate($request, $rules, $message);
        $enterpriseId = $this->enterpriseId;
        $storeId      = $request->post('store_id',0);
        $startTime    = $request->post('start_time', '');
        $endTime      = $request->post('end_time', '');
        $type         = $request->post('type', 1);

        /*
        // 如果是店长则
        if ($memberInfo['is_boss'] == 1) {
            $workerIds = [];
        } else if ($type == 1) {
            $storeId = 0;
        }
        */

        $memberInfo = $this->storeWorkerDetail($enterpriseId, $storeId);

        $selectType=0;
        $workerIds=[];
        // 状态判断
        switch ($this->role) {
            case Constant::LOGIN_ROLE_STORE_MEMBER: // 如果角色是店员
                if ($type == Constant::EXTENSION_ORDER_TYPE) {
                    $storeId = 0;
                }
                $workerIds  = [$memberInfo['id']];
                break;
            case Constant::LOGIN_ROLE_STORE_MANAGER: // 如果角色是店长
                $workerIds=[];
                if ($type == Constant::EXTENSION_ORDER_TYPE) {
                    $selectType = 1;
                }
                break;
            case Constant::LOGIN_ROLE_ENTERPRISE: // 如果企业导购角色
                if ($type == Constant::EXTENSION_ORDER_TYPE) {
                    $storeId = 0;
                }
                $workerIds  = [$memberInfo['id']];
                break;
            case Constant::LOGIN_ROLE_COMMUNITY: // 社群导购角色
                if ($type == Constant::EXTENSION_ORDER_TYPE) {
                    $storeId = 0;
                }
                $workerIds  = [$memberInfo['id']];
                break;
            case Constant::LOGIN_ROLE_ENTERPRISE_ADMIN: // 企业管理角色
                $workerIds = [];
                $storeId   = 0;
                break;
            default:
        }


        $params=[
            'enterprise_id' => $enterpriseId,
            'store_id'      => $storeId,
            'worker_ids'    => $workerIds,
            'start_time'    => $startTime,
            'end_time'      => $endTime,
            'select_type'   => $selectType,
            'type'          => $type
        ];
        Log::info('gw-assistant统计数量统计接口'.json_encode($params));

        $returnData = $storePerformance->storeStatistics($params);
        Log::info('gw-assistant统计数量统计接口返回'.json_encode($returnData));
        return success($returnData['data']);
    }
}
