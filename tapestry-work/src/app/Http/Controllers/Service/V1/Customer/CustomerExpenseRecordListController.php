<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/17
 * Time: 15:40
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Models\CrsOrder;
use CrsSdk\Shop\Order;
use App\Common\Constant;
use Illuminate\Http\Request;
use App\Services\CustomerService;
use App\Services\StaffManagerService;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\Service\V1\BaseController;

class CustomerExpenseRecordListController extends BaseController
{
    /**
     * 客户消费记录列表
     *
     * @param Request $request
     *
     * @return array|mixed
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/3 18:44
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        if (!$enterpriseId) {
            error(600000);
        }
        $hash = $request->input('qy_hash', '');
        if (empty($hash)) {
            error(110002);
        }

        $page = $request->input('page', Constant::DEFAULT_PAGE);
        $size = $request->input('size', Constant::DEFAULT_SIZE);

        $unionId = $request->input('union_id', '');
        if (empty($unionId)) {
            return [
                'list'  => [],
                'page'  => $page,
                'size'  => $size,
                'total' => 0,
            ];
        }

        switchAppHost(['enterprise_id' => $enterpriseId]);

        $list = $this->onlineOrderList($enterpriseId, $unionId, $page, $size);

        return $list;

    }

    /**
     * 线上订单列表
     *
     * @param int    $enterpriseId
     * @param string $unionId
     * @param int    $page
     * @param int    $size
     *
     * @return array|mixed
     * @throws \Exception
     */
    protected function onlineOrderList($enterpriseId = 0, $unionId = '', $page = 0, $size = 0)
    {
        Log::info('获取线上订单列表', [$enterpriseId, $unionId, $page, $size]);
        // 查询客户的userId
       /* if (Constant::HAO_WU_JI_ENTERPRISE == $enterpriseId) {

            $list = CrsOrder::query()
                ->where("sync_source", CrsOrder::HAO_WU_JI_ORDER)
                ->where('union_id', $unionId)
                ->where('pay_at','>',0)
                ->orderBy("updated_at", 'desc')
                ->paginate($size, ['order_no'], '', $page)
                ->toArray();

            if (empty($list['data'])) {
                return [
                    'list'  => [],
                    'page'  => $page,
                    'size'  => $size,
                    'total' => 0,
                ];
            }
            $search = [
                'order_no' => array_column($list['data'], 'order_no')
            ];
        } else {
       */
            $customerInfo = CustomerService::getCustomerListFromShop($enterpriseId, [$unionId]);
            if (empty($customerInfo)) {
                return [
                    'list'  => [],
                    'page'  => $page,
                    'size'  => $size,
                    'total' => 0,
                ];
            }
            $customerUserId = $customerInfo[0]['id'] ?? 0;
            $search         = [
                'enterprise_id' => $enterpriseId,
                'user_id'       => $customerUserId
            ];
      /*  }
*/
        // 查询订单列表
        $orderList = (new Order())->orderList([
            'page'   => $page,
            'size'   => $size,
            'search' => $search
        ]);
        $orderList = $this->setStaffName($enterpriseId, $orderList);
        /*
        // 好物集
        if (Constant::HAO_WU_JI_ENTERPRISE == $enterpriseId) {
            $orderList['total'] = $list['total'] ?? 0;
        }
        */

        $orderList = $this->formatOrder($orderList);

        return $orderList;
    }

    /**
     * setStaffName
     *
     * @desc 设置订单服务人员名称
     *
     * @param $enterpriseId
     * @param $orderList
     *
     * @return mixed
     * @date 2020-07-20
     * @throws \Exception
     */
    private function setStaffName($enterpriseId, $orderList)
    {
        $staffIdList = array_column($orderList['list'], 'service_staff_id');
        $staffIdList = array_unique($staffIdList);
        foreach ($staffIdList as $key => $value) {
            if ($value == 0) {
                unset($staffIdList[$key]);
            }
        }
        $staffIdList = array_values($staffIdList);
        // 从gw_shop查询员工信息
        $staffList = StaffManagerService::getStoreWorkerByIds($enterpriseId, $staffIdList);
        if (!empty($staffList) && isset($staffList['list'])) {
            $staffList = array_column($staffList['list'], null, 'id');
        }
        foreach ($orderList['list'] as $key => $value) {
            $orderList['list'][$key]['waiter_name']   = $staffList[$value['service_staff_id']]['name'] ?? '';
            $orderList['list'][$key]['waiter_mobile'] = $staffList[$value['service_staff_id']]['mobile'] ?? '';
        }
        return $orderList;
    }

    /**
     * formatOrder
     *
     * @desc 格式化订单数据
     *
     * @param $orderList
     *
     * @return mixed
     * @date 2020-07-23
     */
    private function formatOrder($orderList)
    {
        foreach ($orderList['list'] as $key => &$value) {
            $value['create_at']         = $value['created_at'];
            $value['goods_list']        = $value['goods_data'];
            $value['actual_pay_amount'] = $value['cash_fee'] - $value['total_refund_fee'];
            unset($value['goods_data']);
            $value['goods_total'] = $value['quantity'];
            if (empty($value['order_remark'])) {
                $value['order_remark'] = ['content' => ''];
            } else {
                $value['order_remark'] = $value['order_remark'][0];
            }
        }
        return $orderList;
    }

}