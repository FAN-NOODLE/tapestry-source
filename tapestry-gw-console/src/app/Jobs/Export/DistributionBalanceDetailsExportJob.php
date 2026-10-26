<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use CrsSdk\Work\DistributionBalance;
use Illuminate\Support\Facades\Log;

class DistributionBalanceDetailsExportJob extends ExportJob
{

    /**
     * 中心商场无门店版本
     */
    CONST CENTER_VERSION = 3;

    private $version;

    /**
     * GoodsExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     * @throws \Exception
     */
    public function __construct(array $params, array $paramJson = [], string $name = '分销余额管理导购详情导出')
    {
        parent::__construct($params, $paramJson, $name);
        $this->version = $params['version'];
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
     * User: zhaoyang@vchangyi.com
     * Date: 2020/10/17 10:02
     */
    protected function export($params)
    {
        $header = [
            '结算单号',
            '结算状态',
            '订单编号',
            '订单状态',
            '订单支付时间',
            '推广佣金（元）'
        ];
        if($this->version != self::CENTER_VERSION) {
            array_push($header, '服务佣金（元）');
        }
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
     * User: zhaoyang@vchangyi.com
     * Date: 2020/10/16 20:13
     */
    private function listExport($params)
    {
        // 查询获取按导购统计订单列表
        try {
            $result = (new DistributionBalance())->balanceDetailsExport($params);
        } catch (\Exception $exception) {
            Log::info('获取全员分销余额管理导购详情失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取全员分销余额管理导购详情失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $value) {
                $tmpArr = [
                    "\t" . (string)$value['id'],
                    $this->settlementStatus($value['settlement_status']),
                    "\t" . (string)$value['order_no'],
                    $this->orderStatus($value['order_status']),
                    $this->time($value['order_pay_time']),
                    $value['promotion_commission'] / 100
                ];
                if($this->version != self::CENTER_VERSION) {
                    array_push($tmpArr, $value['service_commission'] / 100);
                }
                $cellData[] = $tmpArr;
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

    /**
     * 结算状态维护
     *
     * @param $status
     *
     * @return string
     */
    private function settlementStatus($status)
    {
        $statusName = '待结算';
        switch ($status) {
            case 1:
                $statusName = '已结算';
                break;
            case 2:
                $statusName = '已取消';
                break;
        }
        return $statusName;
    }

    /**
     * 订单状态维护
     *
     * @param $status
     *
     * @return string
     */
    private function orderStatus($status)
    {
        $orderStatusName = [
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
        return $orderStatusName[$status];
    }

    private function time($time)
    {
        return date("Y-m-d H:i:s", $time);
    }

}