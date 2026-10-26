<?php


namespace App\Http\Controllers\App\V1\EnterprisePayment;


use App\Http\Controllers\App\V1\BaseController;
use App\Models\EnterprisePayment;
use App\Services\EnterprisePaymentService\EnterprisePaymentService;
use Maatwebsite\Excel\Facades\Excel;

class ExportController extends BaseController
{
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required | integer'
        ],[
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
        ]);

        $params = $this->request->only([
            'payment_status',
            'enterprise_id',
            'start_time',
            'end_time',
            'time_type',
            'page',
            'limit'
        ]);
        $page = $params['page'] ?? 1;
        $perPage = $params['limit'] ?? 10;
        $list = (new EnterprisePaymentService)->getList($params, $page, $perPage);
        if (empty($list)) {
            return notice(400105);
        }
        return $this->export($list['data']);
    }

    protected function export($data)
    {
        // 表头
        $cellData = [
            ['支付来源', '支付金额(元)', '平台收款金额(元)', '服务费(元)', '明细', '支付方式', '支付流水号', '支付时间', '到账时间', '货款状态', '备注'],
        ];
        $colWidth = [
            'A' => 12, 'B' => 12, 'C' => 15, 'D' => 12, 'E' => 50, 'F' => 15, 'G' => 20, 'H' => 20,
            'I' => 20, 'J' => 15, 'K' => 50
        ];
        // 列表数据
        foreach ($data as $value) {
            switch ($value['payment_category']) {
                case EnterprisePayment::PAYMENT_CATEGORY_RED_PACK:
                    $info = '购买红包可用金额' . number_format(round($value['available_amount'] / 100, 2), 2, '.', '');
                    break;
                case EnterprisePayment::PAYMENT_CATEGORY_SMS:
                    $info = '购买短信可用条数' . $value['buy_number'];
                    break;
                case EnterprisePayment::PAYMENT_CATEGORY_ZHDG:
                    $info = '购买智慧导购账号' . $value['buy_number'] . '个，有效时间' . date('Y-m-d', $value['use_time']) . ' - ' . date('Y-m-d', $value['fail_time']);
                    break;
                case EnterprisePayment::PAYMENT_CATEGORY_SCRM:
                    $info = '购买SCRM账号' . $value['buy_number'] . '个，有效时间' . date('Y-m-d', $value['use_time']) . ' - ' . date('Y-m-d', $value['fail_time']);
                    break;
                case EnterprisePayment::PAYMENT_CATEGORY_MALL:
                    $info = '购买商城账号' . $value['buy_number'] . '个，有效时间' . date('Y-m-d', $value['use_time']) . ' - ' . date('Y-m-d', $value['fail_time']);
                    break;
                case EnterprisePayment::PAYMENT_CATEGORY_SYSTEM:
                    $info = '购买综合系统账号' . $value['buy_number'] . '个，有效时间' . date('Y-m-d', $value['use_time']) . ' - ' . date('Y-m-d', $value['fail_time']);
                    break;
                default:
                    $info = '';
                    break;
            }
            $cellData[] = [
                $value['payment_source'],
                round($value['payment_amount'] / 100, 2),
                round($value['collection_amount'] / 100, 2),
                round($value['server_charge'] / 100, 2),
                $info,
                EnterprisePayment::PAYMENT_TYPE_ARR[$value['payment_type']],
                $value['payment_order_sn'],
                date('Y-m-d H:i:s', $value['payment_time']),
                empty($value['arrival_time']) ? '' : date('Y-m-d H:i:s', $value['arrival_time']),
                ($value['payment_status'] == EnterprisePayment::PAYMENT_STATUS_NO ? '待收款' : '已收款'),
                $value['remark'],
            ];
        }

        $exportFileName = '企业付费列表-' . date('YmdHi', time());
        return Excel::create($exportFileName, function ($excel) use ($cellData, $colWidth) {
            $excel->sheet('score', function ($sheet) use ($cellData, $colWidth) {
                $sheet->setWidth($colWidth)->rows($cellData);

                $endLine = count($cellData);
                // 设置表头样式
                $sheet->cells('A1:V1', function ($cells) {
                    $cells->setAlignment('center')->setValignment('center')->setFontWeight('bold')->setFontSize(13);;
                })->cells('A2:V' . $endLine, function ($cells) { // 设置内容样式
                    $cells->setAlignment('center')->setValignment('center');
                });
            });
        })->export('xlsx');
    }
}