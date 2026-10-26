<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use App\Common\Constant;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use Maatwebsite\Excel\Facades\Excel;

class ExportController extends BaseController
{
    /**
     * 1：订单商品，2：订单赠品
     */
    const ORDER_GOODS_TYPE = [
        1 => '商品',
        2 => '赠品'
    ];

    /**
     * 配送方式：1、快递，2、自提
     */
    const ORDER_DELIVERY_TYPE = [
        1 => '快递',
        2 => '自提'
    ];

    /**
     * 设置导出表头
     *
     * @var array
     */
    protected $colTitle = [
        '门店名称', '门店编码', '订单编号', '主订单编号', '订单状态', '下单账户', '下单时间', '商品编号', '商品名称', '商品规格',
        '商品单价', '商品数量', '小计', '商品总额', '运费', '实付金额拆分总计', '订单总额', '实付金额', '平台实收金额', '支付时间',
        '支付流水号（微信支付）', '收货人姓名', '收货人联系电话', '省', '市', '区', '详细地址', '配送方式', '物流公司', '物流编号',
        '发货时间', '收货时间', '退款总金额', '退款到账时间', '退款流水单号（单号皆对应订单信息非具体商品）', '商品备注'
    ];

    /**
     * 设置导出文件标题
     *
     * @var string
     */
    protected $exportFileName = '订单导出信息列表';

    /**
     * 设置导出列宽
     *
     * @var array
     */
    protected $colWidth = [
        'A'  => 30, 'B' => 30, 'C' => 25, 'D' => 25, 'E' => 12, 'F' => 15, 'G' => 20, 'H' => 20, 'I' => 20,
        'J'  => 15, 'K' => 13, 'L' => 13, 'M' => 12, 'N' => 13, 'O' => 10, 'P' => 13, 'Q' => 20, 'R' => 13,
        'S'  => 18, 'T' => 20, 'U' => 30, 'V' => 15, 'W' => 20, 'X' => 10, 'Y' => 10, 'Z' => 10, 'AA' => 30,
        'AB' => 15, 'AC' => 15, 'AD' => 25, 'AE' => 20, 'AF' => 20, 'AG' => 15, 'AH' => 22, 'AI' => 60, 'AJ' => 13
    ];

    /**
     * 导出筛选的订单
     *
     * @param Request $request
     * @throws \Exception
     */
    public function index(Request $request)
    {
        if ($request->filled('keyword')) {
            $input['search']['keyword'] = $request->get('keyword');
        }

        if ($request->filled('start_time')) {
            $input['search']['start_time'] = $request->get('start_time');
        }

        if ($request->filled('end_time')) {
            $input['search']['end_time'] = $request->get('end_time');
        }

        if ($request->filled('status')) {
            $input['search']['status'] = $request->get('status');
        }

        // 订单编号
        if ($request->filled('order_no')) {
            $input['search']['order_no'] = $request->get('order_no', '');
        }

        // 注册手机号
        if ($request->filled('mobile')) {
            $input['search']['mobile'] = $request->get('mobile', '');
        }

        // 收货人手机号
        if ($request->filled('delivery_mobile')) {
            $input['search']['delivery_mobile'] = $request->get('delivery_mobile', '');
        }

        // 下单人名称
        if ($request->filled('order_user')) {
            $input['search']['order_user'] = $request->get('order_user', '');
        }

        // 下单人名称
        if ($request->filled('delivery_user')) {
            $input['search']['delivery_user'] = $request->get('delivery_user', '');
        }

        // 下单时间排序
        if ($request->filled('order_time_sort')) {
            $input['search']['order_time_sort'] = $request->get('order_time_sort', 0);
        }

        // 支付时间排序
        if ($request->filled('pay_time_sort')) {
            $input['search']['pay_time_sort'] = $request->get('pay_time_sort', 0);
        }

        // SpuId搜索
        if ($request->filled('spu_id')) {
            $input['search']['spu_id'] = $request->get('spu_id', '');
        }

        // SkuId搜索
        if ($request->filled('sku_id')) {
            $input['search']['sku_id'] = $request->get('sku_id', '');
        }

        // 微信优惠券批次号搜索
        if ($request->filled('coupon_stock_id')) {
            $input['search']['coupon_stock_id'] = $request->get('coupon_stock_id', '');
        }

        // 主订单编号
        if ($request->filled('main_order_no')) {
            $input['search']['main_order_no'] = $request->get('main_order_no', '');
        }

        // 配送方式
        if ($request->filled('delivery_type')) {
            $input['search']['delivery_type'] = $request->get('delivery_type', 0);
        }

        // 推广人id
        if ($request->filled('store_staff_id')) {
            $input['search']['store_staff_id'] = $request->get('store_staff_id', 0);
        }

        // 服务人id
        if ($request->filled('service_staff_id')) {
            $input['search']['service_staff_id'] = $request->get('service_staff_id', 0);
        }

        // 订单归属门店id
        if ($request->filled('store_id')) {
            $input['search']['store_id'] = $request->get('store_id', 0);
        }

        // 订单服务门店id
        if ($request->filled('service_store_id')) {
            $input['search']['service_store_id'] = $request->get('service_store_id', 0);
        }

        $input['search']['enterprise_id'] = $this->enterpriseId;
        $input['search']['is_payment'] = Constant::ORDER_PAYMENT_TRUE;
        $input['search']['is_refund'] = Constant::ORDER_REFUND_TRUE;
        $input['search']['is_gift'] = Constant::ORDER_GIFT_TRUE;
        $input['search']['domain'] = $this->qyDomain;
        $input['search']['identifier'] = $this->appIdentifier;
        $input['store'] = $this->getStores($this->enterpriseAdmin['role'], false) ?? [];
        $input['size'] = 0;

        // 订单列表
        $list = (new Order())->orderList($input);
        // 格式化数据
        list($list, $mergeData, $endLine) = $this->formatData($list['list']);

        $this->exports($list, $mergeData, $endLine);
    }

    /**
     * 格式化数据
     *
     * @param array $list
     *
     * @return array
     */
    private function formatData($list = [])
    {
        // 格式化数据、合并订单数组、合并退款数组
        $resultData = $mergeData = [];

        $start = 2;
        $end = 1;
        if ($list) {
            foreach ($list as $key => $val) {
                // 订单状态
                $type = $this->statusConversion((int)$val['status'], $val['refund_at']);

                // 商品总额
                $totalSellPrice = 0;
                // 商品拆分总金额
                $goodsPayTotalPrice = 0;
                foreach ($val['goods_data'] as &$goods) {
                    // 商品
                    if (1 == $goods['goods_type']) {
                        $totalSellPrice += $goods['sell_price'] * $goods['quantity'];
                    }

                    if (!empty($goods['goods_pay_price'])) {
                        $goods['goods_pay_price'] = sprintf("%.2f", $goods['goods_pay_price'] / Constant::MONEY_UNIT);
                    } else {
                        $goods['goods_pay_price'] = 0;
                    }
                    $goodsPayTotalPrice += $goods['goods_pay_price'];
                }

                // 退款总金额
                $refundTotal = array_sum(array_column($val['refund'], 'refund_fee')) ?? 0;

                $isGive = 0;
                //检查订单商品均摊是否都是0
                $check = true;
                foreach ($val['goods_data'] as $item) {
                    if ($item['apportion_price'] != 0) {
                        $check = false;
                        break;
                    }
                }

                $payAt = $val['pay_at'] ?? 0;
                if (empty($val['goods_data'])) {
                    $end++;
                }

                foreach ($val['goods_data'] as $k => $v) {
                    $apportionAmount = '';
                    if (!empty($payAt)) {
                        $apportionAmount = ($v['apportion_price'] * $v['quantity'] + $v['difference_price']) / Constant::MONEY_UNIT;
                    }
                    if ($apportionAmount == 0 && $isGive == 0 && $check) {
                        $apportionAmount = $val['cash_fee'] / Constant::MONEY_UNIT;
                        $isGive = 1;
                    }

                    // 平台实收金额，TODO 临时解决方案，后续需要将主单的other_fee拆分到各子单
                    if ($val['type'] == 16) { // 好物集
                        // 好物集要求导出的平台实收金额不减退款金额
                        $orderPayAmount = ($val['cash_fee'] + $val['other_fee'] - $refundTotal) / Constant::MONEY_UNIT;
                    } else {
                        $orderPayAmount = ($val['cash_fee'] - $refundTotal) / Constant::MONEY_UNIT;
                    }

                    $end++;
                    $resultData[] = [
                        // 门店名称 A
                        $val['store_name'],
                        // 门店编号 B
                        $val['store_no'],
                        // 订单编号 C
                        "\t" . (string)$val['order_no'],
                        // 主订单编号 D
                        "\t" . (string)$val['main_order_no'],
                        // 订单状态 E
                        $type,
                        // 下单账户 F
                        $this->trimStr($val['nick_name']),
                        // 下单时间 G
                        empty($val['created_at']) ? '' : date('Y-m-d H:i:s', (int)$val['created_at']),
                        // 商品编号 H
                        "\t" . (string)$v['sku_no'],
                        // 商品名称 I
                        $v['spu_name'],
                        // 商品规格 J
                        implode(' ', $v['sku_attribute']) ?? '',
                        // 商品单价 K
                        $v['sell_price'] / Constant::MONEY_UNIT,
                        // 商品数量 L
                        $v['quantity'],
                        // 小计 M
                        (int)($v['sell_price'] * $v['quantity']) / Constant::MONEY_UNIT,
                        // 商品总额 N
                        $totalSellPrice / Constant::MONEY_UNIT,
                        // 运费 O
                        $val['logistics_price'] / Constant::MONEY_UNIT,
                        // 实付金额拆分 P
                        $apportionAmount,
                        // 订单总额 Q
                        $val['order_total_amount'] / Constant::MONEY_UNIT,
                        // 实付金额 R
                        $val['cash_fee'] / Constant::MONEY_UNIT,
                        // 平台实收金额 S
                        // ($val['cash_fee'] - $refundTotal) / Constant::MONEY_UNIT,
                        empty($val['pay_at']) ? 0 : $orderPayAmount,
                        // 支付时间 T
                        empty($val['pay_at']) ? '' : date('Y-m-d H:i:s', (int)$val['pay_at']),
                        // 支付流水号（微信支付） U
                        "\t" . (string)$val['third_payment_no'],
                        // 收货人姓名 V
                        !empty($val['delivery_name']) ? $val['delivery_name'] : $val['take_goods_name'] ?? '',
                        // 收货人联系电话 W
                        !empty($val['delivery_mobile']) ? $val['delivery_mobile'] : $val['take_goods_mobile'] ?? '',
                        // 省  X
                        $val['province'],
                        // 市 Y
                        $val['city'],
                        // 区 Z
                        $val['area'],
                        // 详细地址 AA
                        $val['delivery_address'],
                        // 配送方式 AB
                        self::ORDER_DELIVERY_TYPE[$val['delivery_type']] ?? '',
                        // 物流公司 AC
                        $val['order_logistic'][0]['name'] ?? '',
                        // 物流编号 AD
                        "\t" . (string)($val['order_logistic'][0]['logistic_no'] ?? ''),
                        // 发货时间 AE
                        empty($val['send_at']) ? '' : date('Y-m-d H:i:s', (int)$val['send_at']),
                        // 收货时间 AF
                        empty($val['sign_at']) ? '' : date('Y-m-d H:i:s', (int)$val['sign_at']),
                        // 退款总金额 AG
                        ($refundTotal / Constant::MONEY_UNIT) ?? '',
                        // 退款到账时间 AH
                        empty($val['refund'][$k]['end_time']) ? '' : date('Y-m-d H:i:s', (int)$val['refund'][$k]['end_time']),
                        // 退款流水单号 AI
                        "\t" . (string)($val['refund'][$k]['refund_id'] ?? ''),
                        // 商品类型 AJ
                        self::ORDER_GOODS_TYPE[$v['goods_type']] ?? ''
                    ];
                }

                $mergeData[] = [$start, $end];
                $start = $end + 1;
            }
        }

        return [$resultData, $mergeData, $end];
    }

    /**
     * @param array $resultList
     * @param array $mergeData
     * @param       $endLine
     *
     * @return mixed
     */
    public function exports($resultList = [], $mergeData = [], $endLine)
    {
        array_unshift($resultList, $this->colTitle);

        return Excel::create($this->exportFileName, function ($excel) use ($resultList, $mergeData, $endLine) {
            $excel->sheet('score', function ($sheet) use ($resultList, $mergeData, $endLine) {
                $sheet->setWidth($this->colWidth)
                    ->rows($resultList);

                // 合并单元格
                $sheet
                    ->setMergeColumn([
                        'columns' => ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'N', 'O', 'Q', 'R', 'S', 'T',
                            'U', 'V', 'W', 'X', 'Y', 'Z', 'AA', 'AB', 'AC', 'AD', 'AE', 'AF','AG'],
                        'rows'    => $mergeData
                    ])
                    // 设置表头样式
                    ->cells('A1:AJ1', function ($cells) {
                        $cells->setAlignment('center')->setValignment('center')->setFontWeight('bold')->setFontSize(13);;
                    })
                    // 设置内容样式
                    ->cells('A2:AJ' . $endLine, function ($cells) {
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


    /**
     * 订单状态转换文字
     *
     * @param $data
     * @param $refund_at
     *
     * @return string
     */
    public function statusConversion($data, $refund_at)
    {
        $type = '未知状态';
        if (in_array($data, [1, 20, 23])) {
            $type = '待付款';
        } elseif (in_array($data, [25, 30])) {
            $type = '待确认';
        } elseif (in_array($data, [2, 40])) {
            $type = '待发货';
        } elseif (in_array($data, [3, 21])) {
            $type = '已取消';
        } elseif (in_array($data, [5, 200])) {
            $type = '订单退款中';
        } elseif (in_array($data, [6, 100, 205]) && $refund_at > 0) {
            $type = '退款成功';
        } elseif (in_array($data, [7, 201])) {
            $type = '退款失败';
        } elseif (in_array($data, [8, 60])) {
            $type = '已发货';
        } elseif (in_array($data, [9, 80])) {
            $type = '已签收';
        } elseif (in_array($data, [10])) {
            $type = '已关闭';
        } elseif (in_array($data, [11])) {
            $type = '取消订单中';
        } elseif (in_array($data, [41, 61, 81, 82, 84])) {
            $type = '售后中';
        } elseif (in_array($data, [83])) {
            $type = '部分退款';
        } elseif (in_array($data, [12, 100]) && $refund_at == 0) {
            $type = '已完成';
        }
        return $type;
    }
}
