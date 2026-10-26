<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/17
 * Time: 15:40
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\CustomerService;
use App\Services\StaffManagerService;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;

class CustomerDetailController extends BaseController
{
    /**
     * index
     * @desc 顾客详情接口
     * @param Request $request
     *
     * @return array|void
     * @date 2020-07-17
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->get('enterprise_id', 0);
        if (!$enterpriseId) {
            error(600000);
        }
        $storeId = $request->get('store_id', 0);
        $memberMobile = $request->get('mobile');
        if (!$memberMobile) {
            error(600010);
        }
        $unionId = $request->get('union_id', '');
        $relationId = $request->get('relation_id', 0);

        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', Constant::DEFAULT_SIZE);

        switchAppHost(['enterprise_id' => $enterpriseId]);
        // 获取店员信息
        $memberInfo = StoreWorkerService::detail([
            'enterprise_id' => $enterpriseId,
            'store_id' => $storeId ?? 0,
            'mobile' => $memberMobile,
            'worker_status' => StoreWorker::IN_SERVICE
        ]);
        if (empty($memberInfo)) {
            error(600103);
        }

        // 查询企业信息，获取domain
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        if (empty($enterpriseInfo)) {
            error(110008);
        }

        // 查询顾客基本信息
        $customerInfo = CustomerService::customerInfo($enterpriseId, $memberInfo, $unionId, $enterpriseInfo['domain'], $storeId, $relationId);

        // 顾客在商城有信息则查询顾客订单数据
        $orderList = [
            'list' => [],
        ];
        if (!empty($customerInfo['user_id'])) {
            $orderList = (new Order())->orderList([
                'page' => $page,
                'size' => $size,
                'search' => [
                    'enterprise_id' => $enterpriseId,
                    'user_id' => $customerInfo['user_id']
                ]
            ]);
            $orderList = $this->setStaffName($enterpriseId, $orderList);
            $orderList = $this->formatOrder($orderList);
        }
        return success([
            'customer_info' => $customerInfo,
            'order_list' => $orderList,
            'member_info' => $this->formatMemberInfo($memberInfo)
        ]);
    }

    /**
     * setStaffName
     * @desc 设置订单服务人员名称
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
            $orderList['list'][$key]['waiter_name'] = $staffList[$value['service_staff_id']]['name'] ?? '';
            $orderList['list'][$key]['waiter_mobile'] = $staffList[$value['service_staff_id']]['mobile'] ?? '';
        }
        return $orderList;
    }

    /**
     * formatOrder
     * @desc 格式化订单数据
     * @param $orderList
     *
     * @return mixed
     * @date 2020-07-23
     */
    private function formatOrder($orderList)
    {
        foreach ($orderList['list'] as $key => &$value) {
            $value['create_at'] = $value['created_at'];
            $value['goods_list'] = $value['goods_data'];
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

    /**
     * formatMemberInfo
     * @desc 格式化店员信息
     * @param $memberInfo
     *
     * @return array
     * @date 2020-07-23
     */
    private function formatMemberInfo($memberInfo)
    {
        return [
            'id' => $memberInfo['id'],
            'store_id' => $memberInfo['store_id'],
            'type' => $memberInfo['type'],
            'is_boss' => $memberInfo['is_boss'],
            'mem_uid' => $memberInfo['mem_uid'],
        ];
    }
}