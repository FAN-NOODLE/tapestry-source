<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 18:08
 */

namespace App\Http\Controllers\App\V1\Order;


use App\Common\Constant;
use App\Models\StoreMemberViewTab;
use CrsSdk\Shop\Order;
use CrsSdk\Shop\StorePerformance;
use Illuminate\Support\Facades\Auth;

class OrderCountController extends OrderBaseController
{
    /**
     * 获取不同状态订单的数量
     *
     * @param Order $order
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Exception
     */
    public function index(Order $order)
    {
        $storeId   = $this->request->get('store_id', 0);
        $workerId  = $this->request->get('select_work_id', 0);
        $selectAll = $this->request->get('select_all', 0);

        $enterpriseId = $this->enterpriseId;
        $haoWuJiEnterpriseId = env('HAO_WU_JI_ENTERPRISE_ID');
        // 好物集定制减少查库
        if (intval($this->enterpriseId) == $haoWuJiEnterpriseId) {
            $orderData = [
                'all'                  => 0,
                'toBeDelivered'        => 0,
                'onlyPickedUp'         => 0,
                'afterSaleToBeChecked' => 0,
                'toBeRefunded'         => 0,
                'toBeServed'           => 0,
                'extensionNum'         => 0,
                'myOrderNum'           => 0,
                'serviceNum'           => 0
            ];

            return success($orderData);
        }

        $selectType   = 0;
        $memberInfo   = $this->storeWorkerDetail($enterpriseId, $storeId);
        $memberInfoId = $memberInfo['id'] ?? 0;
        // 全部订单的数量
        $orderCountDefault = ['all' => 0];
        $orderCount        = [];
        try {
            // 如果是门店导购则显示数字
            if (in_array($this->role, [Constant::LOGIN_ROLE_STORE_MEMBER, Constant::LOGIN_ROLE_STORE_MANAGER])
                || !empty($workerId) || !empty($selectAll)
            ) {

                $selectWorkerId = 0;

                if (empty($selectAll)) {
                    $selectWorkerId = !empty($workerId) ? $workerId : $memberInfoId;
                }

                $orderCount = $order->statusCount([
                    'store_id'      => $storeId,
                    'enterprise_id' => $enterpriseId,
                    'worker_id'     => $selectWorkerId
                ]);
            }

            $orderCountDefault['toBeDelivered']        = isset($orderCount['to_be_delivered']) ? intval($orderCount['to_be_delivered']) : 0;
            $orderCountDefault['onlyPickedUp']         = isset($orderCount['to_be_picked']) ? intval($orderCount['to_be_picked']) : 0;
            $orderCountDefault['afterSaleToBeChecked'] = isset($orderCount['after_sale_to_be_checked']) ? intval($orderCount['after_sale_to_be_checked']) : 0;
            $orderCountDefault['toBeRefunded']         = isset($orderCount['after_sale_to_be_refund']) ? intval($orderCount['after_sale_to_be_refund']) : 0;
            $orderCountDefault['toBeServed']           = isset($orderCount['to_be_served']) ? intval($orderCount['to_be_served']) : 0;

        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }

        // 全部待处理的订单数量
        $orderCountDefault['all'] = $orderCountDefault['toBeDelivered'] +
            $orderCountDefault['onlyPickedUp'] +
            $orderCountDefault['afterSaleToBeChecked'] +
            $orderCountDefault['toBeRefunded'];

        // 如果是店长并且没有查询其他导购
        if ($this->role == Constant::LOGIN_ROLE_STORE_MANAGER && empty($workerId) || !empty($selectAll)) {
            $selectType = 1;
        }

        // 获取推广数字
        $orderCountDefault['extensionNum'] = $this->getExtensionNum($enterpriseId, $storeId, !empty($workerId) ? $workerId : $memberInfoId, $selectType);

        // 我的订单数量
        $orderCountDefault['myOrderNum'] = $orderCountDefault['all'];
        // 服务订单数量
        $orderCountDefault['serviceNum'] = $orderCountDefault['all'] + $orderCountDefault['toBeServed'];

        return success($orderCountDefault);
    }

    /**
     * 获取推广数字
     *
     * @param $enterpriseId
     * @param $storeId
     * @param $workerId
     * @param $selectType
     *
     * @return mixed
     * @throws \Exception
     */
    private function getExtensionNum($enterpriseId, $storeId, $workerId, $selectType)
    {

        // 查询用户查看推广记录
        $storeTabInfo = StoreMemberViewTab::query()
            ->where('store_worker_id', $workerId)
            ->where('type', StoreMemberViewTab::EXTENSION_TAB_TYPE)
            ->first();

        $returnData = (new StorePerformance())->homeTab([
            'enterprise_id'       => $enterpriseId,
            'store_id'            => $storeId,
            'worker_id'           => $workerId,
            'select_type'         => $selectType,
            'extension_view_time' => $storeTabInfo->extension_view_time ?? 0
        ]);

        return $returnData['data']['extension_tab_num'];
    }

}
