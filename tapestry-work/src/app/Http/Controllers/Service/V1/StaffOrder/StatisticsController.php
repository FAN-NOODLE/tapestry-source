<?php


namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\StaffOrderRequest;
use App\Models\StaffOrderRelation;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\DB;

class StatisticsController extends BaseController
{
    /**
     * 内购员工统计
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->only(['page', 'size', 'search']);
        $searchData = $input['search'] ?? [];
        $search = [
            'type' => $searchData['type'] ?? 0,
            'start_time' => $searchData['start_time'] ?? '',
            'end_time' => $searchData['end_time'] ?? '',
            'enterprise_id' => $searchData['enterprise_id'] ?? '',
        ];

        list($startTime, $endTime) = (new StaffOrderService())->orderTimeWhere($search);

        // 查询订单列表
        $staffOrder = StaffOrderRelation::query()
            ->when(!empty($startTime), function ($query) use ($startTime) {
                return $query->where('created_at', '>=', $startTime);
            })
            ->when(!empty($endTime), function ($query) use ($endTime) {
                return $query->where('created_at', '<', $endTime);
            })
            ->where('enterprise_id', $search['enterprise_id'])
            ->where('order_no', '<>',1)
            ->where('order_pay_time', '>', 0)
            ->get(['refund_money', 'pay_money'])
            ->toArray();

        // 初始化
        $orderAllMoney = $orderSuccessMoney = $orderSuccessTotal = $orderRefundTotal = $orderRefundMoney = 0;

        // 循环订单数据
        foreach ($staffOrder as $value) {
            // 订单总金额（包含退款）
            $orderAllMoney += $value['pay_money'];
            // 支付成功金额（不包含退款）
            if (empty($value['refund_money'])) {
                // 金额
                $orderSuccessMoney += $value['pay_money'];
                // 数量
                $orderSuccessTotal++;
            } else {
                $orderRefundMoney += $value['refund_money'];
                // 未退款数
                $orderRefundTotal++;
            }

        }

        $data = [
            'order_all_total' => count($staffOrder), // 订单总数
            'order_all_money' => $orderAllMoney, // 订单总金额
            'order_success_total' => $orderSuccessTotal, // 不含退款订单数
            'order_success_money' => $orderSuccessMoney, // 不含退款订单金额
            'order_refund_total' => $orderRefundTotal,  // 退款订单数
            'order_refund_money' => $orderRefundMoney, // 退款订单金额
        ];

        return success(['data' => $data]);
    }

}