<?php
/**
 * Created by PhpStorm.
 * User: zq
 * Date: 19-4-8
 * Time: 下午6:52
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use Maatwebsite\Excel\Facades\Excel;
use Illuminate\Http\Request;
use CrsSdk\Shop\Order;

class ExportAfterSaleListController extends BaseController
{
    /**
     * 设置导出列宽
     *
     * @var array
     */
    protected $colWidth = [
        'A' => 20,
        'B' => 20,
        'C' => 20,
        'D' => 15,
        'E' => 20,
        'F' => 15,
        'G' => 20,
        'H' => 20,
        'I' => 15,
        'J' => 20,
        'K' => 20,
        'L' => 20,
        'M' => 15,
        'N' => 100
    ];

    // 申请售后时状态
    protected $applyStatus = [
        '40'  => '待发货',
        '60'  => '待收货',
        '80'  => '已签收',
        '83'  => '部分退款',
        '100' => '已完成'
    ];

    // 售后类型
    protected $applyType = [
        '1' => '仅退款',
        '2' => '退货退款',
    ];

    // 售后订单状态
    protected $applySaleStatus = [
        '21' => '待审核',
        '22' => '待审核',
        '23' => '待用户寄回',
        '24' => '待仓库收货',
        '25' => '人工中',
        '26' => '同意退款',
        '27' => '已驳回',
        '28' => '已取消',
        '29' => '已驳回'
    ];

    // 订单状态
    protected $orderStatus = [
        /**
         * 门店版订单状态
         */
        20  => '待付款',
        21  => '已取消',
        25  => '待确认',
        40  => '待发货',
        41  => '售后中',
        60  => '待收货',
        80  => '已签收',
        81  => '售后中',
        83  => '部分退款',
        100 => '已完成',
        200 => '退款中',
        201 => '退款失败',
        205 => '已退款',

        /**
         * 老版订单状态
         */
        1   => '待付款',
        2   => '待发货',
        3   => '已取消',
        5   => '退款中',
        7   => '退款失败',
        8   => '待收货',
        9   => '已签收',
        11  => '售后中',
        12  => '已完成',
        23  => '待付款',
        30  => '待确认',
    ];

    public function index(Request $request)
    {
        $search = [];
        if ($request->filled('start')) {
            $search['start'] = $request->get('start', '');
        }
        if ($request->filled('end')) {
            $search['end'] = $request->get('end', '');
        }
        if ($request->filled('order_no')) {
            $search['order_no'] = $request->get('order_no', '');
        }
        if ($request->filled('type')) {
            $search['type'] = $request->get('type', 0);
        }
        if ($request->filled('mobile')) {
            $search['mobile'] = $request->get('mobile', '');
        }
        if ($request->filled('delivery_mobile')) {
            $search['delivery_mobile'] = $request->get('delivery_mobile', '');
        }
        if ($request->filled('order_user')) {
            $search['order_user'] = $request->get('order_user', '');
        }
        if ($request->filled('delivery_name')) {
            $search['delivery_name'] = $request->get('delivery_name', '');
        }
        if ($request->filled('status')) {
            $search['status'] = $request->get('status', '');
        }
        if ($request->filled('current_order_status')) {
            $search['current_order_status'] = $request->get('current_order_status', 0);
        }
        if ($request->filled('refund_apply_type')) {
            $search['refund_apply_type'] = $request->get('refund_apply_type', 0);
        }
        if ($request->filled('apply_time_sort')) {
            $search['apply_time_sort'] = $request->get('apply_time_sort', 0);
        }
        if ($request->filled('main_order_no')) {
            $search['main_order_no'] = $request->get('main_order_no', 0);
        }

        $input['search']        = $search;
        $input['enterprise_id'] = $this->enterpriseId;

        $lists = (new Order())->afterSaleList($input)['list'];

        $cellData[] = [
            '订单编号',
            '主订单编号',
            '下单客户',
            '收货人',
            '收货人手机',
            '售后状态',
            '申请时订单状态',
            '售后申请时间',
            '售后方式',
            '申请退款金额',
            '实际退款金额',
            '退款到账时间',
            '订单状态',
            '退款失败原因',
        ];
        if ($lists) {
            foreach ($lists as $value) {
                // 实际退款金额
                $payRefundFee = $value['pay_refund_fee'] / 100 ?? 0;
                // 非"同意退款"
                if (26 != $value['status']) {
                    $payRefundFee = 0;
                }

                $cellData[] = [
                    "\t" . (string)$value['order_no'],
                    "\t" . (string)$value['main_order_no'],
                    $this->trimStr($value['nick_name']),
                    ($value['delivery_type'] == 1) ? ($value['delivery_name'] ?? '') : ($value['take_goods_name'] ?? ''),
                    ($value['delivery_type'] == 1) ? ($value['delivery_mobile'] ?? '') : ($value['take_goods_mobile'] ?? ''),
                    $this->applySaleStatus[$value['status']] ?? '',
                    $this->applyStatus[$value['current_order_status']] ?? '',
                    empty($value['created_at']) ? '' : date('Y-m-d H:i:s', (int)$value['created_at']),
                    $this->applyType[$value['refund_apply_type']] ?? '',
                    $value['current_refund_fee'] / 100 ?? 0,
                    $payRefundFee,
                    empty($value['pay_refund_time']) ? '' : date('Y-m-d H:i:s', (int)$value['pay_refund_time']),
                    $this->orderStatus[$value['order_status']] ?? '',
                    $value['refund_failed_reason'] ?? ''
                ];
            }
        }
        $name = '售后单列表' . date('YmdHi', time());

        return Excel::create($name, function ($excel) use ($cellData) {
            $excel->sheet('售后单列表', function ($sheet) use ($cellData) {
                $sheet->setWidth($this->colWidth)->rows($cellData)
                    // 设置表头样式
                    ->cells('A' . 1 . ':N' . 1, function ($cells) {
                        $cells->setAlignment('center')->setValignment('center')->setFontWeight('bold')->setFontSize(13);;
                    })
                    // 设置内容样式
                    ->cells('A' . 2 . ':N' . count($cellData), function ($cells) {
                        $cells->setAlignment('center')->setValignment('center');
                    });
            });
        })->export('xlsx');
    }

    /**
     * 去除空格和特殊字符
     *
     * @param $str
     *
     * @return string
     */
    private function trimStr($str)
    {
        $search = array(" ", " ", "\n", "\r", "\t", "\/", "=", "^");
        $replace = array("", "", "", "", "", "");
        return $this->filterEmoji(str_replace($search, $replace, $str));
    }

    /**
     * 过滤emoji表情 -- 暂时取消使用
     *
     * @param $emojiStr
     *
     * @return string|string[]|null
     */
    public function filterEmoji($emojiStr)
    {
        $emojiStr = preg_replace_callback('/./u', function (array $match) {
            return strlen($match[0]) >= 4 ? '' : $match[0];
        }, $emojiStr);
        return $emojiStr;
    }
}
