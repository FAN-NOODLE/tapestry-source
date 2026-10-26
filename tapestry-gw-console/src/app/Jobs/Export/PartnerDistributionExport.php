<?php


namespace App\Jobs\Export;


use App\Exceptions\Error;
use CrsSdk\Work\PpartnerDistribution;
use Illuminate\Support\Facades\Log;

class PartnerDistributionExport extends ExportJob
{


    protected $version;
    // 中心商城无门店版
    const  CENTER_VERSION = 3;

    const  ORDER_STATUS_NAME = [
        '20' => '未支付',
        '21' => '已取消',
        '23' => '待支付',
        '25' => '支付中',
        '30' => '待确认',
        '40' => '待发货',
        '41' => '售后中',
        '60' => '待收货',
        '61' => '售后中',
        '80' => '已签收',
        '81' => '售后中',
        '83' => '部分退款',
        '100' => '已完成',
        '200' => '退款中',
        '201' => '退款失败',
        '205' => '已退款',
    ];

    const SETTLEMENT_STATUS = [
        '0' => '待结算',
        '1' => '已结算',
        '2' => '已取消',
    ];

    /**
     * GoodsExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     * @throws \Exception
     */
    public function __construct(array $params, array $paramJson = [], string $name = '提现审批列表')
    {
        parent::__construct($params, $paramJson, $name);
        $this->version = isset($params['version']) ? $params['version'] : '';

    }

    /**
     * 导出Excel列表handle
     * @return bool
     */
    public function handle()
    {
        $this->export($this->params);

    }

    /**
     * 导出
     *
     * @param $params
     * @return bool
     */
    protected function export($params)
    {
        $header = [
            '结算单号',
            '结算状态',
            '订单编号',
            '订单状态',
            '订单支付时间',
            '分销合伙人',
            '分销合伙人佣金（元）',
            '所属导购',
            '导购佣金（元）',
        ];

        // 查询列表
        try {
            $this->excelExportInit($header)->listExport($params);
        } catch (\Exception $exception) {
            Log::info('导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            return true;
        }

        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * 按循环获取数据并插入excel
     *
     * @param $params
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function listExport($params)
    {
        try {
            $result = (new PpartnerDistribution())->Distributionlist($params);
        } catch (\Exception $exception) {
            Log::info('获取结算记录列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取结算记录列表失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $value) {
                $cellInfo = [
                    "\t" . (string)$value['id'],
                    self::SETTLEMENT_STATUS[$value['settlement_status']] ?? '未知结算状态',
                    "\t" . (string)$value['order_no'],
                    self::ORDER_STATUS_NAME[$value['order_status']] ?? '未知订单状态',
                    date('Y-m-d H:i:s', $value['order_pay_time']),
                    $value['partner_name'],
                    $value['partner_commission'] / 100,
                    $value['worker_name'],
                    $value['promotion_commission'] / 100,
                ];
                $cellData[] = $cellInfo;
            }
            $this->exportExcel($cellData);

        }
        // 判断是否取完数据
        if (ceil($result['total'] / $params['size']) > $params['page']) {
            $params['page']++;
            $this->listExport($params);
        }
        $this->exportClose();
    }

}