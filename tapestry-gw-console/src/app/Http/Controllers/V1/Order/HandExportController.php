<?php
/**
 * 人工导出订单列表
 * User: pw
 * Date: 2019-03-27
 * Time: 15:45
 */

namespace App\Http\Controllers\V1\Order;

use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use Maatwebsite\Excel\Facades\Excel;
use App\Http\Controllers\V1\BaseController;
use App\Models\Setting;

class HandExportController extends BaseController
{

    // 数据类型：1=待发货订单列表(默认) 2=待发货退款申请 3=已发货退款申请
    const TYPE_PAYED = 1;
    const TYPE_CANCEL = 2;
    const TYPE_AFTER = 3;

    /**
     * 导出文件名
     * @var array
     */
    protected $exportName = [
        1 => '1-待发货订单列表',
        2 => '2-待发货退款申请',
        3 => '3-已发货退款申请'
    ];

    /**
     * 导出筛选的订单
     * @param Request $request
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

        $input['search']['enterprise_id'] = $this->enterpriseId;

        $input['export_type'] = $request->get('export_type') ?? self::TYPE_PAYED;

        // 订单列表
        $list = (new Order())->HandOrderList($input);

        $this->export($list['data'], $input['export_type']);
    }

    /**
     * Excel导出入口
     * @param $data
     * @param $exportType
     * @return mixed
     */
    public function export($data, $exportType)
    {
        // 表头
        $cellData[] = ['门店名称', '门店编号', '订单编号', '收货人', '收货人联系电话', '省', '市', '区', '详细地址', '订单原价', '运费',
            '订单优惠金额', '微信实付金额', '订单实付金额', '商品名称', '商品编号','商品规格', '商品数量',
            '商品单价', '商品总价', '下单时间', '付款时间', '支付方式', '订单备注'];

        // 列宽
        $colWidth = [
            'A' => 25, 'B' => 25, 'C' => 20, 'D' => 15, 'E' => 20, 'F' => 15, 'G' => 15, 'H' => 15,
            'I' => 40, 'J' => 13, 'K' => 10, 'L' => 18, 'M' => 25, 'N' => 18, 'O' => 18, 'P' => 30,
            'Q' => 30, 'R' => 15, 'S' => 13, 'T' => 13, 'U' => 20, 'V' => 20,  'W'=> 20, 'X'=>30,
        ];

        foreach ($data as $value) {
            // 列表数据
            $cellData[] = [
                $value['store_name'] ?? '',
                $value['store_no'] ?? '',
                "\t" . (string)$value['order_no'] ?? '',
                $value['delivery_name'] ?? '',
                $value['delivery_mobile'] ?? '',
                $value['province'] ?? '',
                $value['city'] ?? '',
                $value['area'] ?? '',
                $value['delivery_address'] ?? '',
                $value['total_amount'] ? sprintf('%.2f', (float)$value['total_amount'] / 100) : '',
                $value['logistics_price'] ? sprintf('%.2f', (float)$value['logistics_price'] / 100) : '',
                $value['discount_fee'] ? sprintf('%.2f', (float)$value['discount_fee'] / 100) : '',
                $value['cash_fee'] ? sprintf('%.2f', (float)$value['cash_fee'] / 100) : '',
                $value['final_amount'] ? sprintf('%.2f', (float)$value['final_amount'] / 100) : '',
                $value['spu_name'] ?? '',
                "\t" . (string)$value['sku_no'] ?? '',
                $value['sku_attribute'] ? "\t" .(string)implode('  ',$value['sku_attribute']): '',
                $value['quantity'] ?? '',
                $value['original_price'] ? sprintf('%.2f', (float)$value['original_price'] / 100) : '',
                $value['sub_price'] ? sprintf('%.2f', (float)$value['sub_price'] / 100) : '',
                empty($value['created_at']) ? '' : date('Y-m-d H:i:s', (int)$value['created_at']),
                empty($value['pay_at']) ? '' : date('Y-m-d H:i:s', (int)$value['pay_at']),
                $value['pay_type'] ?? '',
                $value['order_remark'] ?? '',
            ];
        }

        $exportName = $this->exportName[$exportType];
        $exportFileName = $exportName . date('YmdHi', time());

        return Excel::create($exportFileName, function ($excel) use ($cellData, $colWidth) {
            $excel->sheet('score', function ($sheet) use ($cellData, $colWidth) {
                $sheet->setWidth($colWidth)->rows($cellData);

                $endLine = count($cellData);
                // 设置表头样式
                $sheet->cells('A1:W1', function ($cells) {
                    $cells->setAlignment('center')->setValignment('center')->setFontWeight('bold')->setFontSize(13);;
                })->cells('A2:W' . $endLine, function ($cells) { // 设置内容样式
                    $cells->setAlignment('center')->setValignment('center');
                });
            });
        })->export('xlsx');
    }
}
