<?php
/**
 * Created by PhpStorm.
 * User: qihongchuan
 * Date: 2020/02/21
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\DoveDream;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\DoveDream;
use Maatwebsite\Excel\Facades\Excel;

class OrderController extends BaseController
{
    // 订单状态 ：1 待付款，2 已付款，3 已取消, 4 正在退款，5 退款成功，6 退款失败，7:已关闭
    const ORDER_STATUS_LIST = [1 => '待付款', 2 => '已付款', 3 => '已取消', 4 => '正在退款', 5 => '退款成功', 6 => '退款失败', 7 => '已关闭'];

    /**
     * 梦工厂订单列表
     * @author qihongchuan
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function list(Request $request)
    {
        $params = $request->only(['type','keyword','start_time','end_time','status','sort','page','size']);
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getOrderList($params);
        return success($result);
    }

    /**
     * 梦工厂订单导出
     * @author qihongchuan
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function export(Request $request)
    {
        $params = $request->only(['type','keyword','start_time','end_time','status']);
        $params['enterprise_id'] = $this->enterpriseId;

        // 查询数据
        $orderList = (new DoveDream)->getOrderExport($params)['data'];

        // 列宽
        $colWidth = ['A' => 25,'B' => 20,'C' => 15,'D' => 20,'E' => 15,'F' => 15,'G' => 15,'H' => 40,'I' => 13,'J' => 10,'K' => 18];

        // 导出数据
        $cellData[] = ['订单编号','下单客户','下单时间','订单总额','实付金额','支付时间','退款总额','退款到账时间','参加人数','状态'];
        // todo : 字段展示/导出格式调整，如金额等
        foreach ($orderList as $order) {
            $cellData[] = [
                "\t".$order['order_no'] ?? '',
                $order['user_name'] ?? '',
                empty($order['created_at']) ? '' : date('Y-m-d H:i:s', (int)$order['created_at']),
                $order['total_amount'] ?? '',
                $order['final_amount'] ?? '',
                empty($order['pay_at']) ? '' : date('Y-m-d H:i:s', (int)$order['pay_at']),
                $order['refund_fee'],
                empty($order['refund_at']) ? '' : date('Y-m-d H:i:s', (int)$order['refund_at']),
                $order['quantity'] ?? '',
                $this::ORDER_STATUS_LIST[$order['status']] ?? '',
            ];
        }
        // TODO : 导出文件名乱码问题
        $exportFileName = 'order-list-export-' . date('YmdHi', time());

        return Excel::create($exportFileName, function ($excel) use ($cellData, $colWidth) {
            $excel->sheet('score', function ($sheet) use ($cellData, $colWidth) {
                $sheet->setWidth($colWidth)->rows($cellData);

                $endLine = count($cellData);
                // 设置表头样式
                $sheet->cells('A1:X1', function ($cells) {
                    $cells->setAlignment('center')->setValignment('center')->setFontWeight('bold')->setFontSize(13);;
                })->cells('A2:X' . $endLine, function ($cells) { // 设置内容样式
                    $cells->setAlignment('center')->setValignment('center');
                });
            });
        })->export('xlsx');
    }

    /**
     * 梦工厂订单详情
     * @author zhaoqi
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $input = $request->only([
            'order_no',
        ]);

        $rules = [
            'order_no' => 'required'
        ];

        $message = [
            'order_no.required'  => '订单号必须',
        ];

        $this->validate($request, $rules, $message);

        $input['enterprise_id'] = $this->enterpriseId;
        $result = (new DoveDream)->getDetail($input);
        return success($result);
    }

    /**
     * 梦工厂订单详情
     * @author zhaoqi
     * @date 2020/02/21 10:04
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function approval(Request $request)
    {
        $input = $request->only([
            'order_no',
            'refund_id',
            'status',
            'reject'
        ]);

        $rules = [
            'order_no' => 'required',
            'refund_id' => 'required|integer',
            'status' => 'required|in:2,3'
        ];

        $message = [
            'order_no.required'  => '订单号必须',
            'refund_id.required'  => '申请id必须',
            'refund_id.integer'  => '申请id类型错误',
            'status.required'  => '状态必须',
            'status.in'  => '状态类型错误',
        ];

        $this->validate($request, $rules, $message);

        $input['enterprise_id'] = $this->enterpriseId;
        $result = (new DoveDream)->approval($input);
        return success([]);
    }

}

