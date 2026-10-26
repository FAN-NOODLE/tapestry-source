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
use App\Models\StoreMemberViewTab;
use CrsSdk\Shop\Order;
use CrsSdk\Shop\StorePerformance;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;


/**
 * Class StatisticsController
 *
 * @package App\Http\Controllers\App\V1\Store
 */
class HomeOrderTabController extends Controller
{
    /**
     * 业绩数据统计
     * @param Request          $request
     * @param StorePerformance $storePerformance
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function __invoke(Request $request, StorePerformance $storePerformance)
    {

        $enterpriseId = $this->enterpriseId;
        $storeId      = $request->post('store_id',0);
        $memberInfo   = $this->storeWorkerDetail($enterpriseId, $storeId);
        $workerId     = $memberInfo['id']??0;

        // 查询用户查看推广记录
        $storeTabInfo = StoreMemberViewTab::query()
            ->where('store_worker_id', $workerId)
            ->where('type', StoreMemberViewTab::EXTENSION_TAB_TYPE)
            ->first();

        // 如果是店长则
      /*  if ($memberInfo['is_boss'] == 1) {
            $extensionStoreId = $storeId;
            $workerId         = 0;
        } else {
            $extensionStoreId = 0;
        }
*/
        $selectType=0;
        $extensionStoreId = $extensionWorkerId=0;
        // 状态判断
        switch ($this->role) {
            case Constant::LOGIN_ROLE_STORE_MEMBER: // 如果角色是店员
                $extensionStoreId = 0;
                $extensionWorkerId = $workerId;
                break;
            case Constant::LOGIN_ROLE_STORE_MANAGER: // 如果角色是店长
                $extensionStoreId = $storeId;
                $selectType = 1;
                $workerId = 0;
                $extensionWorkerId = $workerId;
                break;
            case Constant::LOGIN_ROLE_ENTERPRISE: // 如果企业导购角色
                $extensionStoreId = 0;
                $extensionWorkerId = $workerId;
                break;
            case Constant::LOGIN_ROLE_COMMUNITY: // 社群导购角色
                $extensionStoreId  = 0;
                $extensionWorkerId = $workerId;
                break;
            case Constant::LOGIN_ROLE_ENTERPRISE_ADMIN: // 企业管理角色
                $storeId           = 0;
                $extensionStoreId  = 0;
                $workerId          = 0;
                $extensionWorkerId = 0;
                break;
            default:
        }

        $returnData = $storePerformance->homeTab([
            'enterprise_id'       => $enterpriseId,
            'store_id'            => $extensionStoreId,
            'worker_id'           => $extensionWorkerId,
            'select_type'         => $selectType,
            'extension_view_time' => $storeTabInfo->extension_view_time ?? 0
        ]);

        Log::info('gw-assistant统计服务角标参数'.json_encode([
                'store_id'      => $storeId,
                'enterprise_id' => $enterpriseId,
                'worker_id'     => $workerId
            ]));

        $orderCount = (new Order())->statusCount([
            'store_id'      => $storeId,
            'enterprise_id' => $enterpriseId,
            'worker_id'     => $workerId
        ]);

        Log::info('gw-assistant统计服务角标接口返回值'.json_encode($orderCount));

        // 全部待处理的订单数量
        $orderAllCount = $orderCount['to_be_delivered'] +
            $orderCount['to_be_picked'] +
            $orderCount['after_sale_to_be_checked'] +
            $orderCount['after_sale_to_be_refund'] + $orderCount['to_be_served'];

        $returnData['data']['service_tab_num'] = +$orderAllCount;

        return success($returnData['data']);
    }
}
